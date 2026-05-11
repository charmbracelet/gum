{ pkgs }:

pkgs.buildGoModule rec {
  pname = "gum";
  version = "0.15.2";

  src = ./.;

  vendorHash = "sha256-TK2Fc4bTkiSpyYrg4dJOzamEnii03P7kyHZdah9izqY=";

  ldflags = [ "-s" "-w" "-X=main.Version=${version}" ];
}
