#!/usr/bin/env python3

import os
import re
import sys
from datetime import datetime, timezone

BLOCKS_DIR = os.path.join(os.path.dirname(__file__), "blocks")
HISTORY_DIR = os.path.join(BLOCKS_DIR, "history")
INDEX_PATH = os.path.join(os.path.dirname(__file__), "index.md")
GRAPH_PATH = os.path.join(os.path.dirname(__file__), "index-graph.md")


def parse_frontmatter(path):
    with open(path) as f:
        content = f.read()
    m = re.match(r"^---\s*\n(.*?)\n---", content, re.DOTALL)
    if not m:
        return None
    fm_text = m.group(1)
    fm = {}
    requires = []
    in_requires = False
    for line in fm_text.strip().split("\n"):
        if line.startswith("requires:"):
            in_requires = True
            continue
        if in_requires:
            if line.strip().startswith("- "):
                req = line.strip()[2:].strip()
                if req and req != "[]":
                    requires.append(req)
            elif ":" in line or not line.strip():
                in_requires = False
                if ":" in line:
                    k, v = line.split(":", 1)
                    fm[k.strip()] = v.strip()
            else:
                in_requires = False
        else:
            if ":" in line:
                k, v = line.split(":", 1)
                fm[k.strip()] = v.strip()
    fm["_requires"] = requires
    fm["_path"] = path
    return fm


def classify_status(raw):
    try:
        return ("chain", int(raw))
    except (ValueError, TypeError):
        return ("web", raw)


def load_blocks():
    blocks = {}
    for base in [BLOCKS_DIR, HISTORY_DIR]:
        if not os.path.isdir(base):
            continue
        for fname in os.listdir(base):
            if fname.endswith(".md") and fname != "index.md":
                path = os.path.join(base, fname)
                fm = parse_frontmatter(path)
                if fm and "id" in fm:
                    blocks[fm["id"]] = fm
    return blocks


def validate(blocks):
    errors = []
    ids = set(blocks.keys())
    chain_nums = []

    for bid, b in blocks.items():
        if bid != os.path.basename(b["_path"]).replace(".md", ""):
            errors.append(f"ID mismatch: {bid} != {os.path.basename(b['_path'])}")
        for req in b["_requires"]:
            if req not in ids:
                errors.append(f"{bid}: requires '{req}' — not found")
            if req == bid:
                errors.append(f"{bid}: requires itself")

        raw = b.get("status", "draft")
        kind, val = classify_status(raw)
        if kind == "chain":
            if isinstance(val, int) and val >= 1:
                chain_nums.append(val)
            else:
                errors.append(f"{bid}: chain status must be int >= 1, got {raw}")

    if chain_nums:
        chain_nums.sort()
        if chain_nums[0] != 1:
            errors.append(f"chain starts at {chain_nums[0]}, must start at 1")
        for i in range(1, len(chain_nums)):
            if chain_nums[i] != chain_nums[i-1] + 1:
                errors.append(f"chain gap: {chain_nums[i-1]} → {chain_nums[i]}")
        dups = [n for n in chain_nums if chain_nums.count(n) > 1]
        if dups:
            errors.append(f"duplicate chain statuses: {list(set(dups))}")

    dup_names = {}
    for bid in blocks:
        name = bid.rsplit(".", 1)[0]
        dup_names.setdefault(name, []).append(bid)
    for name, bids in dup_names.items():
        if len(bids) > 1:
            errors.append(f"duplicate name '{name}': {bids}")

    return errors


def check_cycles(blocks):
    visited = set()
    stack = set()
    cycles = []

    def dfs(bid):
        visited.add(bid)
        stack.add(bid)
        for req in blocks.get(bid, {}).get("_requires", []):
            if req not in visited:
                dfs(req)
            elif req in stack:
                cycles.append(list(stack) + [req])
        stack.discard(bid)

    for bid in blocks:
        if bid not in visited:
            dfs(bid)
    return cycles


def is_available(block, blocks):
    kind, _ = classify_status(block.get("status"))
    if kind == "chain":
        return False
    if block.get("status") not in ("planned",):
        return False
    for req in block["_requires"]:
        rb = blocks.get(req)
        if not rb:
            return False
        rkind, _ = classify_status(rb.get("status"))
        if rkind != "chain":
            return False
    return True


def generate_report(blocks):
    available = []
    in_progress = []
    blocked = []
    chain = []
    web = []
    rejected = []

    for bid, b in sorted(blocks.items()):
        raw = b.get("status", "draft")
        kind, val = classify_status(raw)
        if kind == "chain":
            chain.append((val, bid))
        elif raw == "rejected":
            rejected.append(bid)
        elif raw == "in-progress":
            in_progress.append(bid)
        elif is_available(b, blocks):
            available.append(bid)
        elif raw in ("planned",):
            blocked.append(bid)
        else:
            web.append(bid)

    chain.sort(key=lambda x: x[0])
    chain_ids = [bid for _, bid in chain]
    return available, in_progress, blocked, chain, rejected, web


def write_index(blocks, available, in_progress, blocked, chain, rejected, cycles):
    now = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
    lines = ["# Blockwebchain Index", "", f"Generated: {now}", ""]

    if cycles:
        lines.append("## CYCLES DETECTED")
        for c in cycles:
            lines.append(f"  {' → '.join(c)} → ...")
        lines.append("")

    max_chain = chain[-1][0] if chain else 0
    next_code = max_chain + 1
    lines.append(f"Chain: 1–{max_chain} | Next status code: **{next_code}**")
    lines.append("")

    lines.append("## CHAIN (History)")
    if chain:
        for n, bid in chain:
            b = blocks[bid]
            lines.append(f"- `{n:02d}` **{bid}** — {b.get('title', bid)}")
    else:
        lines.append("(empty)")
    lines.append("")

    lines.append("## AVAILABLE NEXT")
    if available:
        for bid in available:
            b = blocks[bid]
            lines.append(f"- **{bid}** — {b.get('title', bid)}")
    else:
        lines.append("(none)")
    lines.append("")

    if in_progress:
        lines.append("## IN PROGRESS")
        for bid in in_progress:
            b = blocks[bid]
            reqs = b.get("_requires", [])
            req_str = ", ".join(reqs) if reqs else "none"
            lines.append(f"- **{bid}** — {b.get('title', bid)}  (requires: {req_str})")
        lines.append("")

    if blocked:
        lines.append("## BLOCKED")
        for bid in blocked:
            b = blocks[bid]
            missing = [r for r in b.get("_requires", [])
                       if classify_status(blocks.get(r, {}).get("status"))[0] != "chain"]
            lines.append(f"- **{bid}** — {b.get('title', bid)}  (missing: {', '.join(missing)})")
        lines.append("")

    if rejected:
        lines.append("## REJECTED")
        for bid in rejected:
            b = blocks[bid]
            lines.append(f"- **{bid}** — {b.get('title', bid)}")
        lines.append("")

    chain_count = len(chain)
    total_web = len(available) + len(in_progress) + len(blocked) + len(rejected)
    lines.append(f"chain: {chain_count} | web: {total_web} | next code: {next_code}")

    with open(INDEX_PATH, "w") as f:
        f.write("\n".join(lines) + "\n")


def write_graph(blocks, available, in_progress, blocked, chain, rejected, cycles):
    lines = ["```mermaid", "graph TD"]
    safe = {}
    for bid in blocks:
        safe[bid] = bid.replace(".", "_").replace("-", "_")

    chain_ids = [bid for _, bid in chain]

    for n, bid in chain:
        b = blocks[bid]
        label = f"{n:02d} {b.get('title', bid)[:40]}"
        lines.append(f"    {safe[bid]}[\"{label}\"]")
        lines.append(f"    style {safe[bid]} fill:#d4edda,stroke:#28a745")

    for bid in available:
        b = blocks[bid]
        label = b.get('title', bid)[:40]
        lines.append(f"    {safe[bid]}[\"{label}\"]")
        lines.append(f"    style {safe[bid]} fill:#fff3cd,stroke:#ffc107")

    for bid in in_progress:
        b = blocks[bid]
        label = b.get('title', bid)[:40]
        lines.append(f"    {safe[bid]}[\"{label}\"]")
        lines.append(f"    style {safe[bid]} fill:#cce5ff,stroke:#0d6efd")

    for bid in blocked:
        b = blocks[bid]
        label = b.get('title', bid)[:40]
        lines.append(f"    {safe[bid]}[\"{label}\"]")
        lines.append(f"    style {safe[bid]} fill:#f8d7da,stroke:#dc3545")

    lines.append("")

    for bid, b in blocks.items():
        for req in b.get("_requires", []):
            if req in safe:
                lines.append(f"    {safe[req]} --> {safe[bid]}")

    lines.append("")
    for n, bid in enumerate(chain_ids):
        if n + 1 < len(chain_ids):
            nxt = chain_ids[n + 1]
            lines.append(f"    {safe[bid]} -.->|chain| {safe[nxt]}")

    with open(GRAPH_PATH, "w") as f:
        f.write("\n".join(lines) + "\n```\n")


def main():
    blocks = load_blocks()
    if not blocks:
        print("No blocks found.")
        return

    errors = validate(blocks)
    if errors:
        print("VALIDATION ERRORS:")
        for e in errors:
            print(f"  {e}")
        print()

    cycles = check_cycles(blocks)
    available, in_progress, blocked, chain, rejected, web = generate_report(blocks)

    print("BLOCKWEBCHAIN")
    print("─" * 40)

    max_chain = chain[-1][0] if chain else 0
    print(f"Chain: 1–{max_chain}  |  Next code: {max_chain + 1}")

    if cycles:
        print("\nCYCLES:")
        for c in cycles:
            print(f"  {' → '.join(c)}")

    print("\nCHAIN (History):")
    if chain:
        for n, bid in chain:
            b = blocks[bid]
            print(f"  {n:02d}  {bid}  — {b.get('title', bid)}")
    else:
        print("  (empty)")

    print("\nAVAILABLE NEXT:")
    if available:
        for bid in available:
            b = blocks[bid]
            print(f"  {bid}  — {b.get('title', bid)}")
    else:
        print("  (none)")

    if in_progress:
        print("\nIN PROGRESS:")
        for bid in in_progress:
            b = blocks[bid]
            print(f"  {bid}  — {b.get('title', bid)}")

    if blocked:
        print("\nBLOCKED:")
        for bid in blocked:
            b = blocks[bid]
            missing = [r for r in b.get("_requires", [])
                       if classify_status(blocks.get(r, {}).get("status"))[0] != "chain"]
            print(f"  {bid}  — {b.get('title', bid)}  (missing: {', '.join(missing)})")

    if errors:
        print(f"\nerrors: {len(errors)} — fix before proceeding")

    write_index(blocks, available, in_progress, blocked, chain, rejected, cycles)
    write_graph(blocks, available, in_progress, blocked, chain, rejected, cycles)


if __name__ == "__main__":
    main()
