FROM ubuntu:20.04

# Install dependencies
RUN apt-get update && apt-get install -y \
    wget \
    libgtk-3-0 \
    libnotify4 \
    libnss3 \
    libxss1 \
    libxtst6 \
    xdg-utils \
    libatspi2.0-0 \
    libuuid1 \
    libappindicator3-1 \
    libsecret-1-0 \
    libasound2 \
    xvfb

# Download and install Postman
RUN wget https://dl.pstmn.io/download/latest/linux64 -O postman.tar.gz
RUN tar -xzf postman.tar.gz -C /opt
RUN rm postman.tar.gz
RUN ln -s /opt/Postman/Postman /usr/bin/postman

# Set up a user to run Postman
RUN useradd -m postman
USER postman
WORKDIR /home/postman

# Command to run Postman
CMD ["postman"]