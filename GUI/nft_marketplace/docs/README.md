# NFT Marketplace Documentation

This directory contains design notes and operational guides for the
Synnergy NFT marketplace GUI.

## Configuration

Runtime configuration values are defined via environment variables. See
`../.env.example` for the required variables.

## Architecture

The application is structured using a simple service layer under `src/` and
leverages Redux-style state management for future scalability.
