tag=$2
if [[ -z "tag" ]]; then
  echo "usage: $0 <package-name> <tag>"
  exit 1
fi

package=$1
if [[ -z "$package" ]]; then
  echo "usage: $0 <package-name> <tag>"
  exit 1
fi
package_split=(${package//\// })
package_name=${package_split[-1]}

platforms=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/386"
    "linux/amd64"
    "linux/arm"
    "linux/arm64"
    "linux/riscv64"
    "linux/s390x"
    "windows/amd64"
    "windows/arm64"
)

for platform in "${platforms[@]}"
do
  echo $platform
	platform_split=(${platform//\// })
	GOOS=${platform_split[0]}
	GOARCH=${platform_split[1]}
	output_name=release/builds/$package_name'-'$GOOS'-'$GOARCH
	if [ $GOOS = "windows" ]; then
		output_name+='.exe'
	fi

	env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o $output_name $package
	if [ $? -ne 0 ]; then
   		echo 'An error has occurred. Aborting the script execution...'
		exit 1
	fi
done

for file in ./release/builds/*
do
	mkdir -p ./release/$tag
	fileName="$(basename ${file})"

	if [[ $fileName == *".exe"* ]]; then
		fileName=${fileName/%.exe}
	fi

	zip -9jpr ./release/$tag/$fileName.zip $file ./LICENSE ./config.json.example
done

rm -r ./release/builds
