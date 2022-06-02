# Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

import pyttsx3 # pyttsx3 is a text-to-speech conversion library in Python
import speech_recognition as s #Google Speech API in Python

#Functional programming Model

def text_to_speech(text):
    #engine connects us to hardware in this case 
    eng= pyttsx3.init()
    #Engine created 
    eng.say(text)
    #Runs for small duration of time ohterwise we may not be able to hear
    eng.runAndWait()

    
def speech_to_text():
    r=s.Recognizer()# an object r which recognises the voice
    with s.Microphone() as source:
        #when using with statement. The with statement itself ensures proper acquisition and release of resources
        print(r.recognize_google(audio))
        text_to_speech(r.recognize_google(audio)) 
        
speech_to_text()