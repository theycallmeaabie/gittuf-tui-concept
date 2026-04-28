## What I added

- **Loading screen.** Repo load is async, spinner paints right away.
- **Demo mode.** No repo in cwd -> falls back to hardcoded sample data, footer says so.
- **Status bar.** Breadcrumb on the left, `EDIT` / `READ-ONLY` badge on the right.
- **Trusted Keys view (Trust).** Root key card + signers list with `VALID` / `EXPIRED` badges, enter for detail. Existing TUI never surfaced key expiry.
- **Branch Protection view (Policy).** Card per protected branch: pattern, threshold, approvers, force-push status. Easier to skim than the rule list.
- **Verify Ref.** New main-menu item, mocks `gittuf verify-ref` with pass/fail headline + per-check breakdown. No equivalent existed.

## Menu

```
Main
├── Policy
│   ├── View Rules              (real)
│   └── View Branch Protection  (concept)
├── Trust
│   ├── View Global Rules       (real)
│   └── View Trusted Keys       (concept)
└── Verify Ref                  (concept)
```

`esc` backs out one level.

## Run

```
./gittuf tui            # demo mode if cwd isn't a gittuf repo
gittuf tui              # real submenus work, concept ones still reachable
```

No signing key in git config -> mutation keys (`a` / `e` / `d`) disabled, `READ-ONLY` badge.
