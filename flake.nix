{
    outputs = {self, nixpkgs, ...}:
    let
        system="x86_64-linux";
        pkgs=import nixpkgs {inherit system;};
    in
    {
        packages.${system}.default = pkgs.buildGoModule {
            name = "gum";
            src = self;
            vendorSha256="vvNoO5eABGVwvAzK33uPelmo3BKxfqiYgEXZI7kgeSo=";
        };
        
    };

}
