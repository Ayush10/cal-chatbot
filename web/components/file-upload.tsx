"use client"

import type React from "react"

import { useRef } from "react"
import { Button } from "@/components/ui/button"
import { Paperclip, X, File } from "lucide-react"

interface FileUploadProps {
  onFileSelect: (file: File | null) => void
  selectedFile: File | null
}

export function FileUpload({ onFileSelect, selectedFile }: FileUploadProps) {
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleButtonClick = () => {
    fileInputRef.current?.click()
  }

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] || null
    onFileSelect(file)
  }

  const handleRemoveFile = () => {
    onFileSelect(null)
    if (fileInputRef.current) {
      fileInputRef.current.value = ""
    }
  }

  return (
    <div className="relative">
      <input
        type="file"
        ref={fileInputRef}
        onChange={handleFileChange}
        className="hidden"
        accept="image/*,.pdf,.doc,.docx,.txt"
      />

      {selectedFile ? (
        <div className="flex items-center gap-2 p-2 bg-muted rounded-md">
          <File className="h-4 w-4" />
          <span className="text-xs truncate max-w-[150px]">{selectedFile.name}</span>
          <Button variant="ghost" size="icon" className="h-5 w-5 rounded-full" onClick={handleRemoveFile}>
            <X className="h-3 w-3" />
          </Button>
        </div>
      ) : (
        <Button type="button" variant="ghost" size="icon" onClick={handleButtonClick} className="h-8 w-8">
          <Paperclip className="h-4 w-4" />
          <span className="sr-only">Attach file</span>
        </Button>
      )}
    </div>
  )
}
