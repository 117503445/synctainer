import './App.css'

import { useState } from 'react'
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
import { client } from "twirpscript";
import { PostTask, ReqPostTask, RespPostTask } from "./rpc/synctainer.pb";


const platforms = [
  "linux/amd64",
  "linux/arm64",
]

function App() {
  const host = "https://synctainer-api.117503445.top"
  client.baseURL = host;

  const [image, setImage] = useState("")
  const [platform, setPlatform] = useState("linux/amd64")

  const [btnSyncDisable, setBtnSyncDisable] = useState(false)

  const inputRef = useRef<HTMLInputElement>(null);

  const handleCopy = () => {
    if (inputRef.current && inputRef.current.value) {
      navigator.clipboard.writeText(inputRef.current.value);
      sendToast('success', `Copied to clipboard: ${inputRef.current.value}`)
    }
  };

  const sendToast = (variant: VariantType, msg: string) => {
    enqueueSnackbar(msg, {
      variant: variant,
      anchorOrigin: {
        vertical: 'top',
        horizontal: 'center',
      },
    });
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

            const MOCK = false
            // const MOCK = true // should be false in production

            let response: Response
            let respPostTask: RespPostTask

            if (!MOCK) {
              try {
                // response = await fetch(`${host}`, {
                //   method: 'POST',
                //   body: JSON.stringify({
                //     image: image,
                //     platform: platform,
                //   }),
                //   headers: {
                //     'Content-Type': 'application/json',
                //   }
                // })
                respPostTask = await PostTask({
                  image: image,
                  platform: platform,
                  registry: "",
                  username: "",
                  password: "",
                })
              } catch (error) {
                sendToast('error', `Trigger Image Sync Failed: ${error}`)
                setBtnSyncDisable(false)
                return
              }
            } else {
              await new Promise(resolve => setTimeout(resolve, 1000))
              response = new Response(JSON.stringify({
                image: "registry.cn-hangzhou.aliyuncs.com/117503445-mirror/sync:docker.io.library.mysql.latest"
              }))
            }


            // let resp = { image: "" }

            // let text = await response.text()
            // console.log(text)

            // try {
            //   resp = JSON.parse(text)
            // } catch (error) {
            //   sendToast('error', `Trigger Image Sync Failed: ${text}`)
            //   setBtnSyncDisable(false)
            //   return
            // }


            // let newImage: string = resp.image
            let newImage: string = image // TODO

            // setOpenTriggerReqSnackbar(false)
            // setOpenTriggerRespSnackbar(true)
            sendToast('success', `Trigger Image Sync Successfully`)
            setBtnSyncDisable(false)
            if (inputRef.current) {
              inputRef.current.value = newImage
            }
          }}>Sync</Button>
        <div>
          <Box display="flex" alignItems="center">
            <TextField inputRef={inputRef} variant="outlined" fullWidth disabled multiline size="small" inputProps={{ style: { fontSize: 14 } }} />
            <IconButton onClick={handleCopy} sx={{ height: '100%', marginLeft: '-0px', "&:focus": { outline: 'none' } }}>
              <FileCopyIcon />
            </IconButton>
          </Box>
        </div>

        <Link href="https://github.com/117503445/synctainer/actions/workflows/copy.yml"
          target="_blank">See Progress</Link>
      </Stack>
    </>
  )
}

export default App
