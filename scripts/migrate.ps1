Get-Content .env | ForEach-Object {
    if ($_ -match '^([^#][^=]+)=(.+)$') {
        set-Item -Path "env:$($Matches[1])" -Value $Matches[2]
    }
} 
$command =$args[0]
$name = $args[1]

switch ($command) {
    "up" {
        migrate -path migrations -database $env:DATABASE_URL up
    }
    "down" {
        $count = if($name) { $name } else { "1"}
        write-host "Rolling back $count migration(s). continue? [y/n]"
        $confirm = Read-host
        if ($confirm -eq "y") {
            migrate -path migrations -database $env:DATABASE_URL down $count
        }
    }
    "create" {
        migrate create -ext sql -dir migrations -seq $name
    }
}