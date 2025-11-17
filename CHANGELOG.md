# Changelog

## 1.0.2 - 2025-11-17


### 🐛 Fixes

- Use 'scratch' distroless base for even more optmized image

## 1.0.1 - 2025-11-17


### 🐛 Fixes

- Split dockerfile into stages so the final one contains just the final binary and no dependencies

## 1.0.0 - 2025-11-17


### 🐛 Fixes

- Improved user-agent

## 0.3.3 - 2025-11-17


### 🐛 Fixes

- Updated path to version file

## 0.3.2 - 2025-11-17


### ⚡ Performance

- Move go code in ./app so the image size is smaller

## 0.3.1 - 2025-11-17


### 🐛 Fixes

- Actually use logged-in .docker for every docker command

## 0.3.0 - 2025-11-17


### 🚀 Features

- Proper docker auth and push

## 0.2.2 - 2025-11-17


### 🐛 Fixes

- Updated cliff config for proper changelog

## 0.2.1 - 2025-11-17


### 🐛 Fixes

- Updated cliff config for proper changelog

## 0.2.0 - 2025-11-17


### 🐛 Fixes

- Improved makefile for releases using `cliff`

## 0.1.0 - 2025-11-17


### 🐛 Fixes

- Improved makefile for releases


### 🚀 Features

- Initial version of go code and dockerfile


