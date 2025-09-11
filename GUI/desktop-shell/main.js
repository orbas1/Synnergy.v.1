const { app, BrowserWindow, Menu } = require('electron');
const path = require('path');
const fs = require('fs');
const { authenticate } = require('./auth');

let sessionToken = null;
const windows = new Set();

function loadModules() {
  const configPath = path.join(__dirname, 'modules.json');
  try {
    return JSON.parse(fs.readFileSync(configPath, 'utf8'));
  } catch (err) {
    console.error('Failed to load modules configuration:', err);
    return [];
  }
}

function createWindow(url) {
  const win = new BrowserWindow({
    width: 1024,
    height: 768,
    webPreferences: {
      contextIsolation: true,
      nodeIntegration: false,
    },
  });

  windows.add(win);
  win.on('closed', () => windows.delete(win));

  win.loadURL(`${url}?session=${sessionToken}`);
}

function buildMenu() {
  const modules = loadModules();

  const items = modules.map(({ name, url }) => ({
    label: name,
    click: () => {
      if (!sessionToken) {
        try {
          sessionToken = authenticate(process.env.DESKTOP_SHELL_PRIVATE_KEY);
        } catch (err) {
          console.error('Authentication failed:', err.message);
          return;
        }
      }
      createWindow(url);
    },
  }));

  const menu = Menu.buildFromTemplate([{ label: 'Modules', submenu: items }]);
  Menu.setApplicationMenu(menu);
}

app.whenReady().then(buildMenu);

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('activate', () => {
  if (windows.size === 0) buildMenu();
});
