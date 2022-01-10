$compiler = "C:\go\bin\go.exe"
$rootPath = $MyInvocation.MyCommand.Definition | split-path -parent
$assetsPath = $rootPath + "\assets"
$binPath = $rootPath + "\bin"
$srcPath = $rootPath + "\main\db-generator"

Set-Location $rootPath

$isExists = Test-Path $binPath

if ($isExists) {
    Remove-Item -Path $binPath -Recurse
}

New-Item -Path $rootPath -Name "bin" -Type Directory
New-Item -Path $binPath -Name "sql" -Type Directory

Copy-Item -Path $($assetsPath + '\sql\*') -Destination ($binPath + '\sql')
Copy-Item -Path $($assetsPath + '\schemes\*') -Destination $binPath

& $compiler 'build', '-o', ($binPath+'\app.exe'), $srcPath  





