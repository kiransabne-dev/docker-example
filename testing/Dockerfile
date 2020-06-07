#dockerfile for building image

FROM ubuntu:18.04
LABEL maintainer="kiransabnedev@gmail.com"

# RUN apt-get update -y
# RUN apt-get upgrade -y
# CREATE APP USER ##

# Create the home directory for the new app user.
# RUN mkdir -p /home/app

# # Create an app user so our program doesn't run as root.
# RUN groupadd -r app &&\
#     useradd -r -g app -d /home/app -s /sbin/nologin -c "Docker image user" app

# # Set the home directory to our app user's home.
# ENV HOME=/home/app
# ENV APP_HOME=/home/app/projects

# ## SETTING UP THE APP ##
# RUN mkdir $APP_HOME
# WORKDIR $APP_HOME

# # ***
# # Do any custom logic needed prior to adding your code here
# # ***

# # Copy in the application code.
# ADD . $APP_HOME

# # Chown all the files to the app user.
# COPY --chown=app:app . $APP_HOME

# Change to the app user.
# USER app
# RUN useradd --create-home appuser
# WORKDIR /home/appuser
# USER appuser

# RUN useradd -ms /bin/bash newuser
# USER newuser
# WORKDIR /home/newuser

# ARG USER_ID
# ARG GROUP_ID

# RUN addgroup --gid $GROUP_ID user
# RUN adduser --disabled-password --gecos '' --uid $USER_ID --gid $GROUP_ID user

# # RUN groupadd --gid 5000 newuser \
# #     && useradd --home-dir /home/newuser --create-home --uid 5000 \
# #         --gid 5000 --shell /bin/sh --skel /dev/null newuser

# USER user

# WORKDIR /home/newuser
RUN useradd -ms /bin/bash user

# USER newuser
# WORKDIR /home/newuser

RUN apt-get update -y
RUN apt-get upgrade -y

RUN apt-get install -y gcc
RUN apt-get install -y g++
# RUN apt-get install -y ruby
RUN apt-get install -y python
RUN apt-get install -y nodejs
RUN apt-get install -y npm
# RUN apt-get install -y golang-go	
RUN apt-get install -y curl
RUN apt-get install -y sudo
RUN apt-get install -y nano

USER user

WORKDIR /home/user

run mkdir projects

ENTRYPOINT ["/bin/bash"]