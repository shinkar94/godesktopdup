# Changelog

## v0.3.1 (uncommitted — under review)

### Added
- `disp.Output.GetDesc` — fills `DXGI_OUTPUT_DESC` (device name, desktop coordinates, rotation); needed for diagnostics on rotated panels. (Accidentally omitted from the v0.3.0 tag.)

## v0.3.0 (uncommitted — under review) (uncommitted — under review)

### Fixed
- **Rotated outputs with the desktop image in system memory** (`MapDesktopSurface` path): the mapped surface is in the panel's NATIVE orientation (e.g. portrait rows on a landscape desktop — Bay Trail tablets), while the code copied it linearly into a landscape-shaped buffer — producing interlaced garbage frames on every rotated output. The system-memory path now transposes the frame per the DXGI rotation value (270/90/180), with the native extent re-sliced correctly (the mapped data has `logicalWidth` rows at `pitch = logicalHeight*4`). Verified on a real Bay Trail tablet (1280×800 landscape desktop, rotation 270): clean stream in both orientations; unchanged on identity-rotation machines.
- **Rotation constants shifted vs DXGI** (`disp.ModeRotation`): `Identity` was 0 and `Rotate270` was 3, while DXGI defines UNSPECIFIED=0, IDENTITY=1, ROTATE90=2, ROTATE180=3, ROTATE270=4. With the shifted values the rotation handling branched on wrong modes. Now aligned with DXGI.
- **windows/386 compilation**: the giant-array pointer-cast idiom (`*[1<<30]uint32`) exceeds the 32-bit address space; replaced with `unsafe.Slice` via a shared `viewU32` helper (identical behavior).

### Added
- GitHub Actions: CI (gofmt/vet/build on windows/amd64 + windows/386) and tag-driven Release workflow.
- NOTE: `disp.Output.GetDesc` was accidentally OMITTED from the v0.3.0 tag (present in the working tree, missed in the transfer) — it ships in v0.3.1.

### Notes
- The lab-validated cursor anti-ghost rework is intentionally NOT part of this release (unverified state) — it will follow in a separate version.

## v0.2.1 (uncommitted — under review)

### Fixed
- **Rotate180 in `copyRotatedFrame`** — the branch was a copy-paste of Rotate90 (transpose instead of a 180° flip). 180° now correctly mirrors the frame on both axes.
- **Runtime rotation changes** — `sc.rotation` used to be fixed on the first frame (`sc.rotation == 0`) and never updated. Now a changed `ModeRotation` invalidates the size caches and forces a full frame copy.

### Changed
- Removed the dead `needsSwizzle` field (written but never read).
