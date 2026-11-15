# Changelog

## v0.2.1 (uncommitted — under review)

### Fixed
- **Rotate180 in `copyRotatedFrame`** — the branch was a copy-paste of Rotate90 (transpose instead of a 180° flip). 180° now correctly mirrors the frame on both axes.
- **Runtime rotation changes** — `sc.rotation` used to be fixed on the first frame (`sc.rotation == 0`) and never updated. Now a changed `ModeRotation` invalidates the size caches and forces a full frame copy.

### Changed
- Removed the dead `needsSwizzle` field (written but never read).
