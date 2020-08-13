#!/usr/bin/sh
echo "Installing dynamited..."

# app directories
DYND_APP="/opt/dynamite/dynamited/bin"
DYND_CONF="/etc/dynamite/dynamited"

# create dirs 
echo "Verifying bin directory exists..."
mkdir -p $DYND_APP
mkdir -p $DYND_CONF

# place files 
echo "Installing conf and binary..."
cp ../pkg/conf/config.yml $DYND_CONF/.
cp ../cmd/dynamited $DYND_APP/.
chmod +x $DYND_APP/dynamited

echo "dynamited updated. Start or restart the dynamited service for changes to take effect."