package cli

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	synn "synnergy"
	ilog "synnergy/internal/log"
)

var dataJSON bool

var (
	feedsStore     = newFeedStore()
	resourcesStore = newResourceStore()
)

func init() {
	dataCmd := &cobra.Command{
		Use:   "data",
		Short: "Manage data feeds and resource catalogs",
	}
	dataCmd.PersistentFlags().BoolVar(&dataJSON, "json", false, "output JSON instead of text")
	dataCmd.AddCommand(newFeedCommand(), newResourceCommand())
	rootCmd.AddCommand(dataCmd)
}

func dataPrint(v interface{}, fallback string) {
	if dataJSON {
		b, err := json.MarshalIndent(v, "", "  ")
		if err == nil {
			fmt.Println(string(b))
			return
		}
	}
	if fallback != "" {
		fmt.Println(fallback)
	}
}

func storageDir() (string, error) {
	if dir := os.Getenv("SYN_DATA_DIR"); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return "", err
		}
		return dir, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		tmp := filepath.Join(os.TempDir(), "synnergy-data")
		if err := os.MkdirAll(tmp, 0o755); err != nil {
			return "", err
		}
		return tmp, nil
	}
	dir := filepath.Join(home, ".synnergy", "data")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

// ---- Feed management ----

type feedStore struct {
	mu     sync.Mutex
	feeds  map[string]*synn.DataFeed
	loaded bool
}

func newFeedStore() *feedStore {
	return &feedStore{feeds: make(map[string]*synn.DataFeed)}
}

func (s *feedStore) loadLocked() error {
	if s.loaded {
		return nil
	}
	dir, err := storageDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, "feeds.json")
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		s.loaded = true
		return nil
	}
	if err != nil {
		return err
	}
	var raw map[string]map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	for id, entries := range raw {
		feed := synn.NewDataFeed(id)
		for key, value := range entries {
			feed.Update(key, value)
		}
		s.feeds[id] = feed
	}
	s.loaded = true
	return nil
}

func (s *feedStore) saveLocked() error {
	dir, err := storageDir()
	if err != nil {
		return err
	}
	snapshot := make(map[string]map[string]string, len(s.feeds))
	for id, feed := range s.feeds {
		snapshot[id] = feed.Snapshot()
	}
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return err
	}
	tmp := filepath.Join(dir, "feeds.json.tmp")
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(dir, "feeds.json"))
}

func (s *feedStore) create(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.loadLocked(); err != nil {
		return err
	}
	if _, exists := s.feeds[id]; exists {
		return fmt.Errorf("feed %s already exists", id)
	}
	s.feeds[id] = synn.NewDataFeed(id)
	if err := s.saveLocked(); err != nil {
		return err
	}
	ilog.Info("cli_data_feed_create", "id", id)
	return nil
}

func (s *feedStore) ensure(id string) (*synn.DataFeed, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.loadLocked(); err != nil {
		return nil, false, err
	}
	feed, ok := s.feeds[id]
	if ok {
		return feed, false, nil
	}
	feed = synn.NewDataFeed(id)
	s.feeds[id] = feed
	if err := s.saveLocked(); err != nil {
		return nil, false, err
	}
	ilog.Info("cli_data_feed_ensure", "id", id, "created", true)
	return feed, true, nil
}

func (s *feedStore) apply(id string, entries map[string]string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.loadLocked(); err != nil {
		return err
	}
	feed, ok := s.feeds[id]
	if !ok {
		feed = synn.NewDataFeed(id)
		s.feeds[id] = feed
	}
	for key, value := range entries {
		feed.Update(key, value)
	}
	if err := s.saveLocked(); err != nil {
		return err
	}
	ilog.Info("cli_data_feed_apply", "id", id, "count", len(entries))
	return nil
}

func (s *feedStore) value(id, key string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.loadLocked(); err != nil {
		return "", false, err
	}
	feed, ok := s.feeds[id]
	if !ok {
		return "", false, nil
	}
	val, ok := feed.Get(key)
	return val, ok, nil
}

func (s *feedStore) deleteKey(id, key string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.loadLocked(); err != nil {
		return false, err
	}
	feed, ok := s.feeds[id]
	if !ok {
		return false, nil
	}
	_, exists := feed.Get(key)
	if !exists {
		return false, nil
	}
	feed.Delete(key)
	if err := s.saveLocked(); err != nil {
		return false, err
	}
	ilog.Info("cli_data_feed_delete_key", "id", id, "key", key)
	return true, nil
}

func (s *feedStore) snapshot(id string) (map[string]string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.loadLocked(); err != nil {
		return nil, false, err
	}
	feed, ok := s.feeds[id]
	if !ok {
		return nil, false, nil
	}
	return feed.Snapshot(), true, nil
}

func (s *feedStore) list() ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.loadLocked(); err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(s.feeds))
	for id := range s.feeds {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids, nil
}

func newFeedCommand() *cobra.Command {
	feedCmd := &cobra.Command{Use: "feed", Short: "Manage structured data feeds"}

	createCmd := &cobra.Command{
		Use:   "create [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Provision a new data feed",
		RunE: func(cmd *cobra.Command, args []string) error {
			return feedsStore.create(args[0])
		},
	}

	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply a JSON manifest of key/value pairs to a feed",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("feed")
			manifest, _ := cmd.Flags().GetString("file")
			if id == "" || manifest == "" {
				return errors.New("--feed and --file are required")
			}
			entries, err := readFeedManifest(manifest)
			if err != nil {
				return err
			}
			if len(entries) == 0 {
				return fmt.Errorf("manifest %s contained no entries", manifest)
			}
			return feedsStore.apply(id, entries)
		},
	}
	applyCmd.Flags().String("feed", "", "target feed identifier")
	applyCmd.Flags().String("file", "", "path to JSON manifest")

	getCmd := &cobra.Command{
		Use:   "get [id] [key]",
		Args:  cobra.ExactArgs(2),
		Short: "Retrieve a value from a feed",
		RunE: func(cmd *cobra.Command, args []string) error {
			val, ok, err := feedsStore.value(args[0], args[1])
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("key %s not found in feed %s", args[1], args[0])
			}
			dataPrint(map[string]string{"id": args[0], "key": args[1], "value": val}, val)
			return nil
		},
	}

	deleteCmd := &cobra.Command{
		Use:   "delete [id] [key]",
		Args:  cobra.ExactArgs(2),
		Short: "Delete a key from a feed",
		RunE: func(cmd *cobra.Command, args []string) error {
			removed, err := feedsStore.deleteKey(args[0], args[1])
			if err != nil {
				return err
			}
			if !removed {
				return fmt.Errorf("key %s not found in feed %s", args[1], args[0])
			}
			dataPrint(map[string]string{"status": "deleted", "id": args[0], "key": args[1]}, "deleted")
			return nil
		},
	}

	snapshotCmd := &cobra.Command{
		Use:   "snapshot [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Render all key/value pairs for a feed",
		RunE: func(cmd *cobra.Command, args []string) error {
			snap, ok, err := feedsStore.snapshot(args[0])
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("feed %s not found", args[0])
			}
			if dataJSON {
				dataPrint(snap, "")
				return nil
			}
			keys := make([]string, 0, len(snap))
			for k := range snap {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Printf("%s=%s\n", k, snap[k])
			}
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all configured feed identifiers",
		RunE: func(cmd *cobra.Command, args []string) error {
			ids, err := feedsStore.list()
			if err != nil {
				return err
			}
			dataPrint(map[string][]string{"feeds": ids}, strings.Join(ids, "\n"))
			return nil
		},
	}

	feedCmd.AddCommand(createCmd, applyCmd, getCmd, deleteCmd, snapshotCmd, listCmd)
	return feedCmd
}

func readFeedManifest(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	entries := map[string]string{}
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err == nil && obj != nil {
		for k, v := range obj {
			if err := insertFeedEntry(entries, k, v); err != nil {
				return nil, err
			}
		}
		if len(entries) > 0 {
			return entries, nil
		}
	}
	var arr []map[string]interface{}
	if err := json.Unmarshal(data, &arr); err == nil {
		for _, item := range arr {
			key, ok := item["key"].(string)
			if !ok || key == "" {
				return nil, fmt.Errorf("manifest entry missing key field")
			}
			val, ok := item["value"]
			if !ok {
				return nil, fmt.Errorf("manifest entry %s missing value field", key)
			}
			if err := insertFeedEntry(entries, key, val); err != nil {
				return nil, err
			}
		}
	}
	return entries, nil
}

func insertFeedEntry(entries map[string]string, key string, value interface{}) error {
	switch v := value.(type) {
	case string:
		entries[key] = v
	case float64:
		entries[key] = strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		entries[key] = strconv.FormatBool(v)
	case nil:
		entries[key] = ""
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		entries[key] = string(b)
	}
	return nil
}

// ---- Resource management ----

type resourceMeta struct {
	Key     string    `json:"key"`
	Size    int       `json:"size"`
	Created time.Time `json:"created_at"`
	Updated time.Time `json:"updated_at"`
	Labels  []string  `json:"labels,omitempty"`
	Source  string    `json:"source,omitempty"`
}

type resourceRecord struct {
	Key     string   `json:"key"`
	Data    string   `json:"data"`
	Labels  []string `json:"labels,omitempty"`
	Source  string   `json:"source,omitempty"`
	Created string   `json:"created_at"`
	Updated string   `json:"updated_at"`
}

type resourceStore struct {
	mu      sync.Mutex
	manager *synn.DataResourceManager
	meta    map[string]resourceMeta
	loaded  bool
}

func newResourceStore() *resourceStore {
	return &resourceStore{manager: synn.NewDataResourceManager(), meta: make(map[string]resourceMeta)}
}

func (s *resourceStore) loadLocked() error {
	if s.loaded {
		return nil
	}
	dir, err := storageDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, "resources.json")
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		s.loaded = true
		return nil
	}
	if err != nil {
		return err
	}
	var records []resourceRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return err
	}
	for _, rec := range records {
		payload, err := base64.StdEncoding.DecodeString(rec.Data)
		if err != nil {
			return err
		}
		s.manager.Put(rec.Key, payload)
		created, _ := time.Parse(time.RFC3339Nano, rec.Created)
		if created.IsZero() {
			created = time.Now().UTC()
		}
		updated, _ := time.Parse(time.RFC3339Nano, rec.Updated)
		if updated.IsZero() {
			updated = created
		}
		s.meta[rec.Key] = resourceMeta{
			Key:     rec.Key,
			Size:    len(payload),
			Created: created,
			Updated: updated,
			Labels:  rec.Labels,
			Source:  rec.Source,
		}
	}
	s.loaded = true
	return nil
}

func (s *resourceStore) saveLocked() error {
	dir, err := storageDir()
	if err != nil {
		return err
	}
	records := make([]resourceRecord, 0, len(s.meta))
	for key, meta := range s.meta {
		data, ok := s.manager.Get(key)
		if !ok {
			continue
		}
		records = append(records, resourceRecord{
			Key:     key,
			Data:    base64.StdEncoding.EncodeToString(data),
			Labels:  meta.Labels,
			Source:  meta.Source,
			Created: meta.Created.Format(time.RFC3339Nano),
			Updated: meta.Updated.Format(time.RFC3339Nano),
		})
	}
	sort.Slice(records, func(i, j int) bool { return records[i].Key < records[j].Key })
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	tmp := filepath.Join(dir, "resources.json.tmp")
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(dir, "resources.json"))
}

func (s *resourceStore) put(key string, data []byte, labels []string, source string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.loadLocked(); err != nil {
		return err
	}
	now := time.Now().UTC()
	meta, ok := s.meta[key]
	if !ok {
		meta = resourceMeta{Key: key, Created: now}
	}
	meta.Size = len(data)
	meta.Updated = now
	if len(labels) > 0 {
		meta.Labels = uniqueStrings(labels)
	}
	if source != "" {
		meta.Source = source
	}
	s.manager.Put(key, data)
	s.meta[key] = meta
	if err := s.saveLocked(); err != nil {
		return err
	}
	ilog.Info("cli_data_resource_put", "key", key, "size", meta.Size)
	return nil
}

func (s *resourceStore) delete(key string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.loadLocked(); err != nil {
		return false, err
	}
	if _, ok := s.meta[key]; !ok {
		return false, nil
	}
	s.manager.Delete(key)
	delete(s.meta, key)
	if err := s.saveLocked(); err != nil {
		return false, err
	}
	ilog.Info("cli_data_resource_delete", "key", key)
	return true, nil
}

func (s *resourceStore) info(key string) ([]byte, resourceMeta, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.loadLocked(); err != nil {
		return nil, resourceMeta{}, false, err
	}
	meta, ok := s.meta[key]
	if !ok {
		return nil, resourceMeta{}, false, nil
	}
	data, ok := s.manager.Get(key)
	if !ok {
		return nil, resourceMeta{}, false, nil
	}
	return data, meta, true, nil
}

func (s *resourceStore) list() ([]resourceMeta, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.loadLocked(); err != nil {
		return nil, err
	}
	metas := make([]resourceMeta, 0, len(s.meta))
	for _, meta := range s.meta {
		metas = append(metas, meta)
	}
	sort.Slice(metas, func(i, j int) bool { return metas[i].Key < metas[j].Key })
	return metas, nil
}

func (s *resourceStore) usage() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.loadLocked(); err != nil {
		return 0, err
	}
	return s.manager.Usage(), nil
}

func (s *resourceStore) prune(allowed map[string]struct{}) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.loadLocked(); err != nil {
		return nil, err
	}
	removed := make([]string, 0)
	for key := range s.meta {
		if _, ok := allowed[key]; ok {
			continue
		}
		s.manager.Delete(key)
		delete(s.meta, key)
		removed = append(removed, key)
	}
	if len(removed) > 0 {
		if err := s.saveLocked(); err != nil {
			return nil, err
		}
		ilog.Info("cli_data_resource_prune", "removed", strings.Join(removed, ","))
	}
	sort.Strings(removed)
	return removed, nil
}

func uniqueStrings(in []string) []string {
	set := make(map[string]struct{}, len(in))
	for _, v := range in {
		if v == "" {
			continue
		}
		set[v] = struct{}{}
	}
	out := make([]string, 0, len(set))
	for v := range set {
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}

func newResourceCommand() *cobra.Command {
	resCmd := &cobra.Command{Use: "resource", Short: "Manage binary resource catalog"}

	putCmd := &cobra.Command{
		Use:   "put",
		Short: "Store or update a resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			key, _ := cmd.Flags().GetString("key")
			if key == "" {
				return errors.New("--key is required")
			}
			filePath, _ := cmd.Flags().GetString("file")
			inline, _ := cmd.Flags().GetString("data")
			b64, _ := cmd.Flags().GetString("base64")
			labels, _ := cmd.Flags().GetStringArray("label")
			source, _ := cmd.Flags().GetString("source")
			var payload []byte
			switch {
			case filePath != "":
				var err error
				if filePath == "-" {
					payload, err = io.ReadAll(os.Stdin)
					if err != nil {
						return err
					}
				} else {
					payload, err = os.ReadFile(filePath)
					if err != nil {
						return err
					}
				}
			case inline != "":
				payload = []byte(inline)
			case b64 != "":
				var err error
				payload, err = base64.StdEncoding.DecodeString(b64)
				if err != nil {
					return err
				}
			default:
				return errors.New("provide --file, --data or --base64")
			}
			return resourcesStore.put(key, payload, labels, source)
		},
	}
	putCmd.Flags().String("key", "", "resource identifier")
	putCmd.Flags().String("file", "", "read bytes from file (use - for stdin)")
	putCmd.Flags().String("data", "", "inline string payload")
	putCmd.Flags().String("base64", "", "base64 encoded payload")
	putCmd.Flags().StringArray("label", nil, "resource labels")
	putCmd.Flags().String("source", "", "source description")

	getCmd := &cobra.Command{
		Use:   "get [key]",
		Args:  cobra.ExactArgs(1),
		Short: "Fetch a resource's contents",
		RunE: func(cmd *cobra.Command, args []string) error {
			outPath, _ := cmd.Flags().GetString("out")
			data, meta, ok, err := resourcesStore.info(args[0])
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("resource %s not found", args[0])
			}
			if outPath != "" {
				if err := os.WriteFile(outPath, data, 0o600); err != nil {
					return err
				}
				dataPrint(map[string]any{"status": "written", "path": outPath, "key": args[0]}, fmt.Sprintf("written %s", outPath))
				return nil
			}
			encoded := base64.StdEncoding.EncodeToString(data)
			payload := map[string]any{
				"key":        meta.Key,
				"data":       encoded,
				"size":       meta.Size,
				"labels":     meta.Labels,
				"source":     meta.Source,
				"created_at": meta.Created.Format(time.RFC3339Nano),
				"updated_at": meta.Updated.Format(time.RFC3339Nano),
			}
			dataPrint(payload, encoded)
			return nil
		},
	}
	getCmd.Flags().String("out", "", "write payload to file")

	deleteCmd := &cobra.Command{
		Use:   "delete [key]",
		Args:  cobra.ExactArgs(1),
		Short: "Delete a resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			removed, err := resourcesStore.delete(args[0])
			if err != nil {
				return err
			}
			if !removed {
				return fmt.Errorf("resource %s not found", args[0])
			}
			dataPrint(map[string]string{"status": "deleted", "key": args[0]}, "deleted")
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List catalogued resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			metas, err := resourcesStore.list()
			if err != nil {
				return err
			}
			if dataJSON {
				type record struct {
					Key     string   `json:"key"`
					Size    int      `json:"size"`
					Labels  []string `json:"labels,omitempty"`
					Source  string   `json:"source,omitempty"`
					Created string   `json:"created_at"`
					Updated string   `json:"updated_at"`
				}
				payload := make([]record, 0, len(metas))
				for _, meta := range metas {
					payload = append(payload, record{
						Key:     meta.Key,
						Size:    meta.Size,
						Labels:  meta.Labels,
						Source:  meta.Source,
						Created: meta.Created.Format(time.RFC3339Nano),
						Updated: meta.Updated.Format(time.RFC3339Nano),
					})
				}
				dataPrint(payload, "")
				return nil
			}
			for _, meta := range metas {
				fmt.Printf("%s (%d bytes)\n", meta.Key, meta.Size)
			}
			return nil
		},
	}

	usageCmd := &cobra.Command{
		Use:   "usage",
		Short: "Show total stored bytes",
		RunE: func(cmd *cobra.Command, args []string) error {
			usage, err := resourcesStore.usage()
			if err != nil {
				return err
			}
			dataPrint(map[string]int64{"usage": usage}, fmt.Sprintf("%d", usage))
			return nil
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info [key]",
		Args:  cobra.ExactArgs(1),
		Short: "Show metadata for a resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, meta, ok, err := resourcesStore.info(args[0])
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("resource %s not found", args[0])
			}
			payload := map[string]any{
				"key":        meta.Key,
				"size":       meta.Size,
				"labels":     meta.Labels,
				"source":     meta.Source,
				"created_at": meta.Created.Format(time.RFC3339Nano),
				"updated_at": meta.Updated.Format(time.RFC3339Nano),
				"data_hash":  hashSummary(data),
			}
			dataPrint(payload, fmt.Sprintf("%s (%d bytes)", meta.Key, meta.Size))
			return nil
		},
	}

	importCmd := &cobra.Command{
		Use:   "import",
		Short: "Import resources from a manifest",
		RunE: func(cmd *cobra.Command, args []string) error {
			manifest, _ := cmd.Flags().GetString("manifest")
			if manifest == "" {
				return errors.New("--manifest is required")
			}
			prune, _ := cmd.Flags().GetBool("prune")
			entries, err := readResourceManifest(manifest)
			if err != nil {
				return err
			}
			allowed := make(map[string]struct{}, len(entries))
			for _, entry := range entries {
				if err := resourcesStore.put(entry.Key, entry.Payload, entry.Labels, entry.Source); err != nil {
					return err
				}
				allowed[entry.Key] = struct{}{}
			}
			if prune {
				removed, err := resourcesStore.prune(allowed)
				if err != nil {
					return err
				}
				dataPrint(map[string]any{"imported": len(entries), "removed": removed}, fmt.Sprintf("imported %d", len(entries)))
				return nil
			}
			dataPrint(map[string]any{"imported": len(entries)}, fmt.Sprintf("imported %d", len(entries)))
			return nil
		},
	}
	importCmd.Flags().String("manifest", "", "JSON manifest describing resources")
	importCmd.Flags().Bool("prune", false, "remove resources not listed in manifest")

	pruneCmd := &cobra.Command{
		Use:   "prune",
		Short: "Remove resources not listed in the provided manifest",
		RunE: func(cmd *cobra.Command, args []string) error {
			manifest, _ := cmd.Flags().GetString("manifest")
			if manifest == "" {
				return errors.New("--manifest is required")
			}
			entries, err := readResourceManifest(manifest)
			if err != nil {
				return err
			}
			allowed := make(map[string]struct{}, len(entries))
			for _, entry := range entries {
				allowed[entry.Key] = struct{}{}
			}
			removed, err := resourcesStore.prune(allowed)
			if err != nil {
				return err
			}
			dataPrint(map[string]any{"removed": removed}, strings.Join(removed, "\n"))
			return nil
		},
	}
	pruneCmd.Flags().String("manifest", "", "JSON manifest with resource keys to retain")

	resCmd.AddCommand(putCmd, getCmd, deleteCmd, listCmd, usageCmd, infoCmd, importCmd, pruneCmd)
	return resCmd
}

type resourceManifestEntry struct {
	Key     string
	Payload []byte
	Labels  []string
	Source  string
}

func readResourceManifest(path string) ([]resourceManifestEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(path)
	var raw []map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	entries := make([]resourceManifestEntry, 0, len(raw))
	for _, item := range raw {
		key, ok := item["key"].(string)
		if !ok || key == "" {
			return nil, fmt.Errorf("manifest entry missing key")
		}
		labels := extractLabels(item)
		source, _ := item["source"].(string)
		var payload []byte
		switch {
		case item["path"] != nil:
			pathVal, _ := item["path"].(string)
			if pathVal == "" {
				return nil, fmt.Errorf("entry %s path empty", key)
			}
			full := pathVal
			if !filepath.IsAbs(full) {
				full = filepath.Join(dir, pathVal)
			}
			b, err := os.ReadFile(full)
			if err != nil {
				return nil, err
			}
			payload = b
		case item["base64"] != nil:
			enc, _ := item["base64"].(string)
			b, err := base64.StdEncoding.DecodeString(enc)
			if err != nil {
				return nil, err
			}
			payload = b
		case item["data"] != nil:
			payload = []byte(fmt.Sprint(item["data"]))
		default:
			return nil, fmt.Errorf("entry %s missing data, base64 or path", key)
		}
		entries = append(entries, resourceManifestEntry{Key: key, Payload: payload, Labels: labels, Source: source})
	}
	return entries, nil
}

func extractLabels(item map[string]interface{}) []string {
	raw, ok := item["labels"].([]interface{})
	if !ok {
		return nil
	}
	labels := make([]string, 0, len(raw))
	for _, v := range raw {
		if s, ok := v.(string); ok && s != "" {
			labels = append(labels, s)
		}
	}
	return uniqueStrings(labels)
}

func hashSummary(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	sum := 0
	for _, b := range data {
		sum = (sum*31 + int(b)) % 0x7fffffff
	}
	return fmt.Sprintf("h%08x", sum)
}

// resetDataOperationsState resets CLI state for tests.
func resetDataOperationsState() {
	feedsStore = newFeedStore()
	resourcesStore = newResourceStore()
}
