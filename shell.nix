{ pkgs ? import <nixpkgs> {} }:

with pkgs; mkShell {
  buildInputs = [ which go electrum ];
}
