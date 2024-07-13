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

const platforms = [
  "linux/amd64",
  "linux/arm64",
]

function App() {

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
    </>
  )
}

export default App
