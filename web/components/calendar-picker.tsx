"use client"

import type React from "react"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Calendar } from "@/components/ui/calendar"
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover"
import { CalendarIcon, Check, X } from "lucide-react"
import { format } from "date-fns"

interface CalendarPickerProps {
  onDateSelect: (date: Date | null) => void
  selectedDate: Date | null
}

export function CalendarPicker({ onDateSelect, selectedDate }: CalendarPickerProps) {
  const [open, setOpen] = useState(false)

  const handleSelect = (date: Date | undefined) => {
    onDateSelect(date || null)
    setOpen(false)
  }

  const handleClear = (e: React.MouseEvent) => {
    e.stopPropagation()
    onDateSelect(null)
  }

  return (
    <div className="relative">
      {selectedDate ? (
        <div className="flex items-center gap-2 p-2 bg-muted rounded-md">
          <CalendarIcon className="h-4 w-4" />
          <span className="text-xs">{format(selectedDate, "PPP")}</span>
          <Button variant="ghost" size="icon" className="h-5 w-5 rounded-full" onClick={handleClear}>
            <X className="h-3 w-3" />
          </Button>
        </div>
      ) : (
        <Popover open={open} onOpenChange={setOpen}>
          <PopoverTrigger asChild>
            <Button variant="ghost" size="icon" className="h-8 w-8">
              <CalendarIcon className="h-4 w-4" />
              <span className="sr-only">Open calendar</span>
            </Button>
          </PopoverTrigger>
          <PopoverContent className="w-auto p-0" align="start">
            <div className="p-3">
              <Calendar mode="single" selected={selectedDate || undefined} onSelect={handleSelect} initialFocus />
              <div className="mt-3 flex items-center justify-end gap-2">
                <Button variant="outline" size="sm" onClick={() => setOpen(false)}>
                  Cancel
                </Button>
                <Button
                  size="sm"
                  onClick={() => {
                    if (selectedDate) setOpen(false)
                  }}
                  disabled={!selectedDate}
                  className="gap-1"
                >
                  <Check className="h-3 w-3" />
                  Select
                </Button>
              </div>
            </div>
          </PopoverContent>
        </Popover>
      )}
    </div>
  )
}
