---
version: patch
---

- Feat: add CLI build-hook mechanism — plugins with dreego-plugin.json run pre-build steps during dreego build/run/dev
- Feat: interactive approval prompt before running plugin build commands — saved to dreego-build.json
- Feat: dreego build --yes flag for CI auto-approval
- Feat: CLI discovers plugins from go.mod requires matching github.com/dreego-stack/* prefix