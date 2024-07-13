import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import Stack from '@mui/material/Stack';
import Box from '@mui/material/Box';
import Autocomplete from '@mui/material/Autocomplete';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';

import React, { useRef } from 'react';
import IconButton from '@mui/material/IconButton';
import FileCopyIcon from '@mui/icons-material/FileCopy';
import Link from '@mui/material/Link';
import Snackbar from '@mui/material/Snackbar';
import MuiAlert from '@mui/material/Alert';

const CopyTextField: React.FC = () => {
  const inputRef = useRef<HTMLInputElement>(null);
  const [openSnackbar, setOpenSnackbar] = useState(false);

  const handleCopy = () => {
    if (inputRef.current) {
      navigator.clipboard.writeText(inputRef.current.value);
      setOpenSnackbar(true);
    }
  };

  const handleCloseSnackbar = () => {
    setOpenSnackbar(false);
  };

  return (
    <div>
    <Box display="flex" alignItems="center">
      <TextField inputRef={inputRef} variant="outlined" fullWidth />
      <IconButton onClick={handleCopy} sx={{ height: '100%', marginLeft: '-40px', "&:focus": { outline: 'none' } }}>
        <FileCopyIcon />
      </IconButton>
    </Box>
    <Snackbar
      anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
      open={openSnackbar}
      autoHideDuration={3000}
      onClose={handleCloseSnackbar}
    >
      <MuiAlert onClose={handleCloseSnackbar} severity="success" sx={{ width: '100%' }}>
        Text copied to clipboard!
      </MuiAlert>
    </Snackbar>
  </div>
  );
};


const top100Films = [
  { title: 'The Shawshank Redemption', year: 1994 },
  { title: 'The Godfather', year: 1972 },
  { title: 'The Godfather: Part II', year: 1974 },
]

const platforms = [
  "linux/amd64",
  "linux/arm64",
]

function App() {
  const [count, setCount] = useState(0)

  return (
    <>
      <Stack spacing={2} sx={{ width: 300 }}>

        <TextField id="outlined-basic" label="Image" variant="outlined" />

        <Autocomplete
          id="platform"
          freeSolo
          // options={top100Films.map((option) => option.title)}
          options={platforms}
          renderInput={(params) => <TextField {...params} label="Platform" />}
        />



        <Button variant="contained">Sync</Button>

        <CopyTextField />
        <Link href="https://github.com/117503445/container-copier/actions/workflows/copy.yml">See Progress</Link>
      </Stack>

      {/* <div>
        <a href="https://vitejs.dev" target="_blank">
          <img src={viteLogo} className="logo" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Vite + React</h1>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p> */}
    </>
  )
}

export default App
