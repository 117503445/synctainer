import './App.css'

import { useState } from 'react'
import React, { useRef } from 'react';

import Stack from '@mui/material/Stack';
import Box from '@mui/material/Box';
import Autocomplete from '@mui/material/Autocomplete';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';
import FileCopyIcon from '@mui/icons-material/FileCopy';
import Link from '@mui/material/Link';
import Snackbar from '@mui/material/Snackbar';
import MuiAlert from '@mui/material/Alert';

type TriggerResp = {
  image: string
}


const platforms = [
  "linux/amd64",
  "linux/arm64",
]

function App() {

  const host = "https://containr-copier-kawdvytifk.cn-hangzhou.fcapp.run"

  const [image, setImage] = useState("")
  const [platform, setPlatform] = useState("linux/amd64")

  const [btnSyncDisable, setBtnSyncDisable] = useState(false)

  const inputRef = useRef<HTMLInputElement>(null);
  const [openCopySnackbar, setOpenCopySnackbar] = useState(false);
  const [openTriggerReqSnackbar, setOpenTriggerReqSnackbar] = useState(false);
  const [openTriggerRespSnackbar, setOpenTriggerRespSnackbar] = useState(false);

  const handleCopy = () => {
    if (inputRef.current) {
      navigator.clipboard.writeText(inputRef.current.value);
      setOpenCopySnackbar(true);
    }
  };

  const handleCloseCopySnackbar = () => {
    setOpenCopySnackbar(false);
  };

  const handleCloseTriggerReqSnackbar = () => {
    setOpenTriggerReqSnackbar(false);
  };

  const handleCloseTriggerRespSnackbar = () => {
    setOpenTriggerRespSnackbar(false);
  };

  return (
    <>
      <Stack spacing={2} sx={{ width: 300 }}>

        <TextField label="Image" variant="outlined" onChange={(e) => setImage(e.target.value)} />

        <Autocomplete
          freeSolo
          options={platforms}
          renderInput={(params) => <TextField {...params} label="Platform" />}
          defaultValue={"linux/amd64"}
          onChange={(e, value) => {
            if (value) {
              setPlatform(value)
            }
          }}
        />



        <Button variant="contained"
          disabled={btnSyncDisable}
          onClick={async () => {
            setBtnSyncDisable(true)
            setOpenTriggerReqSnackbar(true)

            console.info(`Syncing ${image} with ${platform}`)

            let response = await fetch(`${host}`, {
              method: 'POST',
              body: JSON.stringify({
                image: "busybox",
                platform: "linux/amd64",
              }),
              headers: {
                'Content-Type': 'application/json',
              }
            })

            // await new Promise(resolve => setTimeout(resolve, 1000))
            // let resp = {
            // }

            console.log(response.status)

            let resp = await response.json()

            let newImage: string = resp.image

            setOpenTriggerReqSnackbar(false)
            setOpenTriggerRespSnackbar(true)
            setBtnSyncDisable(false)
            if (inputRef.current) {
              inputRef.current.value = newImage
            }
          }}>Sync</Button>

        <Snackbar
          anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
          open={openTriggerReqSnackbar}
          autoHideDuration={3000}
          onClose={handleCloseTriggerReqSnackbar}
        >
          <MuiAlert onClose={handleCloseTriggerReqSnackbar} severity="info" sx={{ width: '100%' }}>
            Syncing busybox with linux/amd64
          </MuiAlert>
        </Snackbar>

        <Snackbar
          anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
          open={openTriggerRespSnackbar}
          autoHideDuration={3000}
          onClose={handleCloseTriggerRespSnackbar}
        >
          <MuiAlert onClose={handleCloseTriggerRespSnackbar} severity="success" sx={{ width: '100%' }}>
            Sync success!
          </MuiAlert>
        </Snackbar>

        <div>
          <Box display="flex" alignItems="center">
            <TextField inputRef={inputRef} variant="outlined" fullWidth disabled />
            <IconButton onClick={handleCopy} sx={{ height: '100%', marginLeft: '-40px', "&:focus": { outline: 'none' } }}>
              <FileCopyIcon />
            </IconButton>
          </Box>
          <Snackbar
            anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
            open={openCopySnackbar}
            autoHideDuration={3000}
            onClose={handleCloseCopySnackbar}
          >
            <MuiAlert onClose={handleCloseCopySnackbar} severity="success" sx={{ width: '100%' }}>
              Text copied to clipboard!
            </MuiAlert>
          </Snackbar>
        </div>

        <Link href="https://github.com/117503445/container-copier/actions/workflows/copy.yml"
          target="_blank">See Progress</Link>
      </Stack>
    </>
  )
}

export default App
