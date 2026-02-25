# Vendor MIBs

This directory holds vendor-specific MIB files that are not available from
IETF, IANA, or CableLabs public sources.

## Adding Vendor MIBs

1. Create a subdirectory named after the vendor (e.g., `cisco/`, `arris/`,
   `harmonic/`).
2. Place `.mib` files inside using the project naming convention:

       MIB-NAME@YYYY-MM-DD.mib

   where `YYYY-MM-DD` is the revision date from the MIB's `LAST-UPDATED`
   clause.

3. Create a symlink from the unversioned name to the latest version:

       ln -s MIB-NAME@YYYY-MM-DD.mib MIB-NAME.mib

4. Commit the files directly -- vendor MIBs are tracked in git because
   they are not downloadable from a public URL.

## Naming Convention

| Component    | Example                        |
|--------------|--------------------------------|
| MIB name     | `ARRIS-CMTS-MIB`              |
| Revision     | `2024-03-15`                   |
| Full name    | `ARRIS-CMTS-MIB@2024-03-15.mib` |
| Symlink      | `ARRIS-CMTS-MIB.mib -> ARRIS-CMTS-MIB@2024-03-15.mib` |

## Notes

- Always verify you have the right to redistribute a vendor MIB before
  committing it. Some vendors restrict redistribution.
- If a vendor publishes MIBs at a stable public URL, consider adding
  support for that source in `tools/mib-downloader/` instead.
