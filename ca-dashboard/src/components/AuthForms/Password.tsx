'use client'

import { Eye, EyeOff } from "lucide-react";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { useState } from "react";

export function PasswordField({
  forgetPass,
  LabelName,
  value,
  onChange,
}: {
  forgetPass: boolean;
  LabelName: string;
  value: string;
  onChange: (value: string) => void;
}) {
    const [show, setShow] = useState(false)
    return (
      <div className="grid gap-3">
        <div className="flex items-center">
          <Label htmlFor="password">{LabelName}</Label>
          {forgetPass && (
            <a
              href="#"
              className="ml-auto text-sm underline-offset-4 hover:underline"
            >
              Forgot your password?
            </a>
          )}
        </div>
  
        <div className="relative">
          <Input
            id="password"
            type={show ? "text" : "password"}
            required
            className="pr-10"
            value={value}
            onChange={(e) => onChange(e.target.value)} // 🔑 calls parent function
          />
          <Button
            type="button"
            variant="ghost"
            size="icon"
            className="absolute right-0 top-0 h-full px-2"
            onClick={() => setShow(!show)}
          >
            {show ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
          </Button>
        </div>
      </div>
    );
  }