import './App.css'

import { useState, useEffect } from 'react'
import { useRef } from 'react';

import Stack from '@mui/material/Stack';
import Autocomplete from '@mui/material/Autocomplete';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';
import Link from '@mui/material/Link';
import { SnackbarProvider, enqueueSnackbar, VariantType } from 'notistack';
import { client, TwirpError } from "twirpscript";
import { GetTask, PostTask } from "./rpc/synctainer.pb";
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

  const [targetImage, setTargetImage] = useState(localStorage.getItem('targetImage') || "")
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
    localStorage.setItem('targetImage', targetImage);
    localStorage.setItem('username', username);
    localStorage.setItem('password', password);
  }, [targetImage, username, password]);

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

        <TextField
          required
          label="Target Image"
          variant="outlined"
          value={targetImage}
          placeholder="registry.cn-hangzhou.aliyuncs.com/117503445-mirror/sync"
          onChange={(e) => setTargetImage(e.target.value)}
          multiline // Enables textarea
          minRows={2} // Optional: minimum number of rows
          maxRows={Infinity} // Optional: allows unlimited expansion
          inputProps={{
            style: {
              whiteSpace: 'pre-wrap', // Preserves whitespace and allows line wrap
              wordWrap: 'break-word', // Allows long words to break and wrap
            },
          }}
        />
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

            // 清除之前的定时器
            timerRef.current.forEach(id => clearTimeout(id));
            timerRef.current = [];
            setTagImage("")
            setHashImage("")
            setGithubActionUrl("")

            setBtnSyncDisable(true)

            console.info(`Syncing ${image} with ${platform}`)

            let taskId = ""

            try {
              const respPostTask = await PostTask({
                image: image,
                platform: platform,
                targetImage: targetImage,
                username: username,
                password: password,
              })
              console.log("respPostTask", respPostTask)
              taskId = respPostTask.id
              setTagImage(respPostTask.tagImage)

              // 设置定时轮询
              const delays = [10, 20, 30, 40, 50, 60, 120, 180, 240, 300];
              console.log("delays", delays)
              delays.forEach(delay => {
                const timerId = window.setTimeout(async () => {
                  console.log("GetTask", "delay", delay)
                  try {
                    const respGetTask = await GetTask({ id: taskId });
                    console.log("respGetTask", respGetTask)
                    if (respGetTask.githubActionUrl) {
                      setGithubActionUrl(respGetTask.githubActionUrl);
                    }
                    if (respGetTask.digest) {
                      // 成功获取digest，清理所有定时器
                      console.log("success get digest, clearing all timers");
                      timerRef.current.forEach(id => clearTimeout(id));
                      timerRef.current = [];


                      setHashImage(respGetTask.digest);
                      sendToast('success', `Sync Task Started`)
                    }
                  } catch (error: unknown) {
                    // 错误处理
                    timerRef.current.forEach(id => clearTimeout(id));
                    timerRef.current = [];
                    if (error instanceof TwirpError) {
                      sendToast('error', `Get Task Failed: ${error.msg}`);
                    } else {
                      sendToast('error', `Get Task Failed: ${error}`);
                    }
                  }
                }, delay * 1000);
                timerRef.current.push(timerId);
              });
            } catch (error: unknown) {
              if (error instanceof TwirpError) {
                sendToast('error', `Trigger Image Sync Failed: ${error.msg}`)
              } else {
                sendToast('error', `Trigger Image Sync Failed: ${error}`)
              }
              setBtnSyncDisable(false)
              return
            }

            sendToast('success', `Sync Task Starting`)
            setBtnSyncDisable(false)
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
