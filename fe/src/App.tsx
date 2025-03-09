import './App.css'

import { useState, useEffect } from 'react'
import { useRef } from 'react';

import Stack from '@mui/material/Stack';
import Box from '@mui/material/Box';
import Autocomplete from '@mui/material/Autocomplete';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';
import FileCopyIcon from '@mui/icons-material/FileCopy';
import Link from '@mui/material/Link';
import { SnackbarProvider, enqueueSnackbar, VariantType } from 'notistack';
import { client, TwirpError } from "twirpscript";
import { GetTask, PostTask, RespPostTask } from "./rpc/synctainer.pb";
import Visibility from '@mui/icons-material/Visibility';
import VisibilityOff from '@mui/icons-material/VisibilityOff';
import CopyableTextField from './components/CopyableTextField';

const platforms = [
  "linux/amd64",
  "linux/arm64",
]

function App() {
  const host = "https://synctainer-api.117503445.top"
  client.baseURL = host;

  const [image, setImage] = useState("")
  const [platform, setPlatform] = useState("linux/amd64")

  const [registry, setRegistry] = useState(localStorage.getItem('registry') || "")
  const [username, setUsername] = useState(localStorage.getItem('username') || "")
  const [password, setPassword] = useState(localStorage.getItem('password') || "")

  const [btnSyncDisable, setBtnSyncDisable] = useState(false)

  const [githubActionUrl, setGithubActionUrl] = useState("")

  const timerRef = useRef<number[]>([]); // 新增定时器引用

  // const inputRef = useRef<HTMLInputElement>(null);

  const [tagImage, setTagImage] = useState("")
  const [hashImage, setHashImage] = useState("")


  useEffect(() => {
    // 监听状态变化并更新 localStorage
    localStorage.setItem('registry', registry);
    localStorage.setItem('username', username);
    localStorage.setItem('password', password);
  }, [registry, username, password]);

  useEffect(() => {
    return () => {
      // 清理所有定时器
      timerRef.current.forEach(id => clearTimeout(id));
    };
  }, []);

  const sendToast = (variant: VariantType, msg: string) => {
    enqueueSnackbar(msg, {
      variant: variant,
      anchorOrigin: {
        vertical: 'top',
        horizontal: 'center',
      },
    });
  };


  const [showPassword, setShowPassword] = useState(false);

  const handleClickShowPassword = () => {
    setShowPassword(!showPassword);
  };

  const handleMouseDownPassword = (event: { preventDefault: () => void; }) => {
    event.preventDefault();
  };

  return (
    <>
      <h1><Link href="https://github.com/117503445/synctainer"
        target="_blank">synctainer</Link></h1>
      <SnackbarProvider maxSnack={3} />
      <Stack spacing={2} sx={{ width: 300 }}>

        <TextField required label="Image" variant="outlined" onChange={(e) => setImage(e.target.value)} />

        <Autocomplete
          freeSolo
          options={platforms}
          renderInput={(params) => <TextField {...params} label="Platform" />}
          defaultValue={"linux/amd64"}
          onChange={(_, value) => {
            if (value) {
              setPlatform(value)
            }
          }}
        />

        <TextField required label="Target Registry" variant="outlined" value={registry} onChange={(e) => setRegistry(e.target.value)} />
        <TextField required label="Username" variant="outlined" value={username} onChange={(e) => setUsername(e.target.value)} />
        <TextField required label="Password" type={showPassword ? 'text' : 'password'} variant="outlined" value={password}
          InputProps={{
            endAdornment: (
              <IconButton
                aria-label="toggle password visibility"
                onClick={handleClickShowPassword}
                onMouseDown={handleMouseDownPassword}
              >
                {showPassword ? <Visibility /> : <VisibilityOff />}
              </IconButton>
            ),
          }}
          onChange={(e) => setPassword(e.target.value)} />

        <Button variant="contained"
          disabled={btnSyncDisable}
          onClick={async () => {

            if (!image) {
              sendToast('error', `Image is required`)
              return
            }

            sendToast('info', `Trigger Image Sync`)

            setBtnSyncDisable(true)

            console.info(`Syncing ${image} with ${platform}`)

            let taskId = ""
            let newImage: string = "" // TODO


            try {
              const respPostTask = await PostTask({
                image: image,
                platform: platform,
                registry: registry,
                username: username,
                password: password,
              })
              console.log("respPostTask", respPostTask)
              taskId = respPostTask.id
              newImage = respPostTask.tagImage

              // 清除之前的定时器
              timerRef.current.forEach(id => clearTimeout(id));
              timerRef.current = [];

              // 设置定时轮询
              const delays = [10, 20, 30, 40, 50, 60, 120, 180, 240, 300];
              delays.forEach(delay => {
                const timerId = window.setTimeout(async () => {
                  try {
                    const respGetTask = await GetTask({ id: taskId });
                    if (respGetTask.digest) {
                      // 成功获取digest，清理所有定时器
                      timerRef.current.forEach(id => clearTimeout(id));
                      timerRef.current = [];
                      setBtnSyncDisable(false);
                      setGithubActionUrl(respGetTask.githubActionUrl);
                      sendToast('success', `Image Sync Completed: ${respGetTask.digest}`);

                      setHashImage(respGetTask.digest);
                    }
                  } catch (error) {
                    // 错误处理
                    timerRef.current.forEach(id => clearTimeout(id));
                    timerRef.current = [];
                    setBtnSyncDisable(false);
                    if (error instanceof TwirpError) {
                      sendToast('error', `Get Task Failed: ${error.msg}`);
                    } else {
                      sendToast('error', `Get Task Failed: ${error}`);
                    }
                  }
                }, delay * 1000);
                timerRef.current.push(timerId);
              });
            } catch (error) {
              if (error instanceof TwirpError) {
                sendToast('error', `Trigger Image Sync Failed: ${error.msg}`)
              } else {
                sendToast('error', `Trigger Image Sync Failed: ${error}`)
              }
              setBtnSyncDisable(false)
              return
            }

            sendToast('success', `Trigger Image Sync Successfully`)
            setBtnSyncDisable(false)
            setTagImage(newImage)
          }}>Sync</Button>

        <CopyableTextField value={tagImage} onChange={(e) => setTagImage(e.target.value)} />
        <CopyableTextField value={hashImage} onChange={(e) => setHashImage(e.target.value)} />

        {githubActionUrl && (
          <Link
            href={githubActionUrl}
            target="_blank"
          >
            See Progress
          </Link>
        )}
      </Stack >
    </>
  )
}

export default App
