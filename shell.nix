{ pkgs ? import <nixpkgs> {} }:

with pkgs; mkShell {
  SSL_CERT_FILE = "${cacert}/etc/ssl/certs/ca-bundle.crt";
  GIT_SSL_CAINFO = "${cacert}/etc/ssl/certs/ca-bundle.crt";

  GO111MODULE = "on";
  GOPATH = "/tmp/lib-cryptocurrency-nix-shell-pure-gopath";

  buildInputs = [ git which go electrum ];
}
