{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  nativeBuildInputs = with pkgs; [
    go protobuf protoc-gen-go protoc-gen-go-grpc grpcurl
  ];
}
