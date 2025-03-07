FROM 117503445/dev-golang

RUN pacman -Syu --noconfirm npm yarn pnpm

RUN npm config set registry https://registry.npmmirror.com
RUN yarn config set registry https://registry.npmmirror.com
RUN pnpm config set registry https://registry.npmmirror.com
# RUN yarn global add typescript

RUN npm install @serverless-devs/s -g

RUN curl -L https://github.com/regclient/regclient/releases/latest/download/regctl-linux-amd64 -o /usr/bin/regctl && chmod +x /usr/bin/regctl
