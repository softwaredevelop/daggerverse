# Get the current working directory
$rootDir = Get-Location

Write-Host "🔍 Searching for 'dagger.json' files..."

# Recursively search for dagger.json files
$daggerFiles = Get-ChildItem -Path $rootDir -Recurse -Filter "dagger.json"

if ($daggerFiles.Count -eq 0) {
    Write-Host "⚠️ No 'dagger.json' files found. Exiting..." -ForegroundColor Yellow
    exit 0
}

# Run in each directory where dagger.json is found
foreach ($file in $daggerFiles) {
    $directory = $file.DirectoryName
    Write-Host "📂 Entering directory: $directory"
    Push-Location $directory

    Write-Host "🚀 Running 'dagger develop'..."
    try {
        dagger develop

        if ($LASTEXITCODE -ne 0) {
            Write-Host "❌ 'dagger develop' failed in $directory" -ForegroundColor Red
        } else {
            Write-Host "✅ 'dagger develop' completed successfully in $directory" -ForegroundColor Green
        }
    } catch {
        Write-Host "🔥 An error occurred while running 'dagger develop' in `${directory}: $_" -ForegroundColor Red
    } finally {
        Pop-Location
    }
}

Write-Host "🎉 All 'dagger develop' executions completed!"
