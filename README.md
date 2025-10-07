# go4crt.sh

`go4crt.sh` is a small, blazing-fast subdomain enumeration utility that pulls certificate transparency records from crt.sh and turns them into a clean, deduplicated list of subdomains—ready for live checking, screenshotting, fuzzing, or any next-step recon. I built it to give pentesters and bug bounty hunters a minimal, reliable tool that does one thing very well: extract subdomains from public TLS certificates and save them to a file so you can plug the results straight into your workflow.

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

## 🎬 Credits 
``` 
Thanks to @TaurusOmar from whom i got an idea and applied it 
:) Thanks https://crt.sh/
Thanks to Yash Pawar for retesting and supporting me 
Thanks to Sujan and jheel for being there :) ````

