{ pkgs }:

pkgs.buildGoModule {
  name = "gum";
  src = ./.;
  vendorSha256 = "sha256-rOBwhPXo4sTSI3j3rn3c5qWGnGFgkpeFUKgtzKBltbg=";
}
