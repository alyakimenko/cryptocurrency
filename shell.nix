# Copyright 2020 Mikhail Klementev. All rights reserved.
# Use of this source code is governed by a MIT license
# that can be found in the LICENSE file.

{ pkgs ? import <nixpkgs> {} }:

with pkgs; mkShell {
  SSL_CERT_FILE = "${cacert}/etc/ssl/certs/ca-bundle.crt";
  GIT_SSL_CAINFO = "${cacert}/etc/ssl/certs/ca-bundle.crt";

  GO111MODULE = "on";
  GOPATH = "/tmp/lib-cryptocurrency-nix-shell-pure-gopath";

  buildInputs = [ git which go electrum ];
}
