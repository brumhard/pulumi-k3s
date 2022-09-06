{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
    # nativeBuildInputs is usually what you want -- tools you need to run
    nativeBuildInputs = with pkgs; [ 
		dotnet-sdk
		pulumictl
		yarn
		nodejs
		python310

		# needed for github.com/mattn/go-ieproxy CGO
		darwin.apple_sdk.frameworks.CFNetwork
		darwin.apple_sdk.frameworks.Security
    ];

	shellHook = ''
    export CGO_ENABLED=1
	'';
}
