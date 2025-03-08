FROM 117503445/dev-golang

RUN pacman -Syu --noconfirm npm yarn pnpm

RUN npm config set registry https://registry.npmmirror.com
RUN yarn config set registry https://registry.npmmirror.com
RUN pnpm config set registry https://registry.npmmirror.com
# RUN yarn global add typescript

RUN npm install @serverless-devs/s -g

RUN curl -L https://github.com/regclient/regclient/releases/latest/download/regctl-linux-amd64 -o /usr/bin/regctl && chmod +x /usr/bin/regctl

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && go install github.com/twitchtv/twirp/protoc-gen-twirp@latest
RUN pacman -Sy --noconfirm protobuf

RUN yarn global add twirpscript

RUN go install github.com/aliyun/ossutil@master