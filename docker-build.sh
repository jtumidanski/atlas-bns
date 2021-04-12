if [[ "$1" = "NO-CACHE" ]]
then
   docker build --no-cache --tag atlas-bns:latest .
else
   docker build --tag atlas-bns:latest .
fi
