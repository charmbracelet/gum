{ pkgs }:

pkgs.buildGoModule {
  name = "gum";
  src = ./.;
  vendorSha256="rMqhYZMa0+5F3X4WDm4jE6IwlzOugqm65SAP38bdQx8=";
}
