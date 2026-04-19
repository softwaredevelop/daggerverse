# List of modules
$modules = @(
    "actionlint",
    "editorconfig",
    "hadolint",
    "hello",
    "quarto",
    "revive",
    "ruff",
    "shellcheck",
    "ssh-manager",
    "yamllint"
)

# Get the current working directory
$rootDir = Get-Location

# Run tests in each module's test directory
foreach ($module in $modules) {
    $testPath = Join-Path $rootDir $module "test"

    if (Test-Path $testPath -PathType Container) {
        Write-Host "📂 Entering test directory: $testPath"
        Set-Location $testPath

        Write-Host "🚀 Running 'dagger call all' in $module/test..."
        dagger call all

        if ($LASTEXITCODE -ne 0) {
            Write-Host "❌ Dagger tests failed in $module/test" -ForegroundColor Red
        } else {
            Write-Host "✅ Dagger tests passed in $module/test" -ForegroundColor Green
        }

        # Return to the root directory
        Set-Location $rootDir
    } else {
        Write-Host "⚠️ Test directory not found: $testPath, skipping..." -ForegroundColor Yellow
    }
}

Write-Host "🎉 All Dagger tests completed!"
