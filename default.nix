{ pkgs }:

pkgs.buildGoModule rec {
  pname = "gum";
  version = "0.15.0";

  src = ./.;

  vendorHash = "sha256-i/KBe41ufYA+tqnB5LCC1geIc2Jnh97JLXcXfBgxdM4=";

  ldflags = [ "-s" "-w" "-X=main.Version=${version}" ];
}
