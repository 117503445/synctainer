FROM 117503445/dev-golang

RUN pacman -Syu --noconfirm npm yarn pnpm

RUN npm config set registry https://registry.npmmirror.com
RUN yarn config set registry https://registry.npmmirror.com
RUN pnpm config set registry https://registry.npmmirror.com
# RUN yarn global add typescript

RUN npm install @serverless-devs/s -g