# go4crt.sh

`go4crt.sh` is a fast and lightweight subdomain finder that extracts data from [crt.sh](https://crt.sh).  
It helps penetration testers and bug bounty hunters enumerate subdomains quickly and save results directly to a file.

---

## ✨ Features
- 🔍 Fetch subdomains from crt.sh Certificate Transparency logs  
- 📂 Save results in a text file  
- ⚡ Simple and fast (pure Go, no dependencies)  

---

## 📦 Installation

```bash
go install github.com/4ncurze/go4crt.sh@latest

```

## 🚀 Usage

```bash
go4crt.sh -d example.com -o subdomains.txt
```
