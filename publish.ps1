git tag | foreach-object -process { git push origin --delete $_ }
git tag | foreach-object -process { git tag -d $_ }
git tag v1.3.0
git push
git push --tags