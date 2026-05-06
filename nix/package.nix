{
  lib,
  buildGoModule,
}:

let
  manifest = lib.importJSON ../.release-please-manifest.json;
  hashes = lib.importJSON ./hashes.json;
in
buildGoModule {
  pname = "fred-mcp";
  version = manifest.".";
  src = ../.;
  vendorHash = hashes.vendorHash;
  subPackages = [ "cmd/fred-mcp" ];

  env.CGO_ENABLED = "0";

  ldflags = [
    "-s"
    "-w"
    "-X main.version=${manifest."."}"
  ];

  meta = with lib; {
    description = "MCP server for Federal Reserve Economic Data (FRED)";
    homepage = "https://github.com/shanehull/fred-mcp";
    license = licenses.mit;
    mainProgram = "fred-mcp";
    platforms = [
      "aarch64-linux"
      "x86_64-linux"
      "aarch64-darwin"
      "x86_64-darwin"
    ];
  };
}
