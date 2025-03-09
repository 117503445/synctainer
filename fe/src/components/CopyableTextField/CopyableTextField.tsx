import React, { useRef } from 'react';
import { TextField, IconButton, Box } from '@mui/material';
import FileCopyIcon from '@mui/icons-material/FileCopy';
import { enqueueSnackbar, VariantType } from 'notistack';

interface CopyableTextFieldProps {
  value?: string;
  onChange?: (event: React.ChangeEvent<HTMLInputElement>) => void;
}

const CopyableTextField: React.FC<CopyableTextFieldProps> = ({ value, onChange }) => {
  const inputRef = useRef<HTMLInputElement>(null);
  const sendToast = (variant: VariantType, msg: string) => {
    enqueueSnackbar(msg, {
      variant: variant,
      anchorOrigin: {
        vertical: 'top',
        horizontal: 'center',
      },
    });
  };
  const handleCopy = () => {
    if (inputRef.current && value) {
      navigator.clipboard.writeText(value).then(() => {
        sendToast('success', `Copied to clipboard: ${value}`)
      }).catch(err => {
        console.error('Failed to copy text: ', err);
      });
    }
  };

  if (!value) return null; // 如果没有内容，则不渲染组件

  return (
    <Box display="flex" alignItems="center">
      <TextField
        inputRef={inputRef}
        variant="outlined"
        fullWidth
        disabled
        multiline
        size="small"
        value={value}
        onChange={onChange}
        inputProps={{ style: { fontSize: 14 } }}
      />
      <IconButton onClick={handleCopy} sx={{ height: '100%', marginLeft: '-0px' }} aria-label="copy">
        <FileCopyIcon />
      </IconButton>
    </Box>
  );
};

export default CopyableTextField;