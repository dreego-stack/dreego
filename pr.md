---
version: patch
---

- Feat: add core KVStore interface (kv-store.1) — `KVStore` interface (Get/Set/Delete/Expire with TTL), like `database/sql`, interface only, plugins implement (Redis/Ristretto/In-Memory), distinct from Storage (blobs)
- Chore: TODO reorg — observability.1 + api-swagger.1 moved to new "Plugins (external repos)" section
