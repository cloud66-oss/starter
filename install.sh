#!/bin/bash
version=1.4.1

# bash utilities credit: http://natelandau.com/bash-scripting-utilities/


#Formatting output
bold=$(tput bold)
reset=$(tput sgr0)
red=$(tput setaf 1)
green=$(tput setaf 76)
tan=$(tput setaf 3)
e_success() { printf "${bold}${green}%s${reset}\n" "$@"
}
e_error() { printf "${bold}${red}%s${reset}\n" "$@"
}
e_warning() { printf "${tan}%s${reset}\n" "$@"
}
e_info() { printf "${bold}%s${reset}\n" "$@"
}

#Check target
type_exists() {
if [ $(type -P $1) ]; then
  return 0
fi
return 1
}
is_os() {
  if [[ "${OSTYPE}" == $1* ]]; then
    return 0
  fi
  return 1
}
is_64() {
  if [ `uname -m` == 'x86_64' ]; then
    return 0 # 64-bit stuff here
  fi
  return 1 # 32-bit stuff here
}

do_install() {
  if is_os "darwin"; then
    e_info "Downloading Starter for macOS..."
    `rm -f /tmp/starter &> /dev/null`
    `curl -L --progress-bar -o /tmp/starter https://s3.amazonaws.com/downloads.cloud66.com/starter/darwin_amd64_v$version`    
  elif is_os "linux"; then
    if [[ is_64 ]]; then
	  e_info "Downloading Starter for Linux x64..."
	  `rm -f /tmp/starter &> /dev/null`
	  `curl -L --progress-bar -o /tmp/starter https://s3.amazonaws.com/downloads.cloud66.com/starter/linux_amd64_v$version`
    else
	  e_error "Aborted: 32 bit version of starter is not currently supported!"
    fi
  else
  	e_error "Aborted: Unable to detect your operating system and architecture!"
  	e_warning "Please download Starter manually from: https://github.com/cloud66-oss/starter/releases"
  	exit 1
  fi
  # extract the archive to local home
  printf "Copying Starter to $USER_HOME/.starter/starter ...\n"
  `mkdir -p $USER_HOME/.starter`
  `rm -f $USER_HOME/.starter/starter &> /dev/null`
  `cp /tmp/starter  $USER_HOME/.starter/starter &> /dev/null`
  printf "Making starter command executable ...\n"
  if [ $UID -eq 0 ] ; then
    `chown $SUDO_USER $USER_HOME/.starter`
    `chown $SUDO_USER $USER_HOME/.starter/starter`
  fi
  `chmod +x $USER_HOME/.starter/starter`
  printf "Creating Starter symlink in /usr/local/bin/starter ...\n"
  `unlink /usr/local/bin/starter &> /dev/null`
  `ln -nfs $USER_HOME/.starter/starter /usr/local/bin/starter &> /dev/null`
  if [ $? -eq 0 ] ; then
  	e_info "The 'starter' command should now be available"
  	e_success "Successfully installed Starter! Go build some images!"
  else
	e_warning "Warning: Unable to create a symlink for Starter"
	e_warning "Please create your symlink manually from $USER_HOME/.starter/starter"
  fi
}

e_success "Installing Starter V$version"
# check if running as sudoer
if [ $UID -eq 0 ] ; then
	USER_HOME="/home/"$SUDO_USER
else
	USER_HOME=$HOME
fi
if type_exists 'tar'; then
  do_install
else
  e_error "Aborted: 'tar' is required to extract the binary. Please install 'tar' first"
  exit 1
fi
printf "\n"
