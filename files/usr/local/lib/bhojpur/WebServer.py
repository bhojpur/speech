#!/usr/bin/python3
# -*- coding: utf-8 -*-

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

# Class WebServer: serves the GUI and process Bhojpur Speech API requests

# --- System-Imports   -------------------------------------------------------

import os, json, queue, traceback, uuid

from flask import Flask, Response, render_template, request, make_response
from flask import send_from_directory, send_file

from werkzeug.serving import make_server, WSGIRequestHandler

from bhojpur import Base

class WebServer(Base):
  """ Serve GUI and process Bhojpur Speech API requests """

  # --- constructor   --------------------------------------------------------

  def __init__(self,app):
    """ constructor """

    self._app          = app
    self._api          = app.api
    self.debug         = app.debug

    self.stop_event    = app.stop_event
    self.read_config(app.options.pgm_dir)
    self._flask = Flask('webspeech',template_folder=self._web_root,
                        root_path=self._web_root)
    self._flask.debug = self.debug
    self._set_routes()

  # --- read configuration   --------------------------------------------------

  def read_config(self,pgm_dir):
    """ read configuration from config-file """

    # section [WEB]
    self._port = int(self.get_value(self._app.parser,"WEB","port",8026))
    self._host = self.get_value(self._app.parser,"WEB","host","0.0.0.0")

    default_web_root = os.path.realpath(
      os.path.join(pgm_dir,"..","lib","bhojpur","web"))
    self._web_root  = self.get_value(self._app.parser,"WEB","web_root",
                                         default_web_root)

  # --- set up routing   -----------------------------------------------------

  def _set_routes(self):
    """ set up routing (decorators don't seem to work for methods within a class """

    # web-interface (GUI)
    self._flask.add_url_rule('/','main',self.main_page)
    self._flask.add_url_rule('/css/<path:filepath>','css',self.css_pages)
    self._flask.add_url_rule('/webfonts/<path:filepath>','webfonts',self.webfonts)
    self._flask.add_url_rule('/images/<path:filepath>','images',self.images)
    self._flask.add_url_rule('/js/<path:filepath>','js',self.js_pages)
    self._flask.add_url_rule('/api/get_events','get_events',self.get_events)
    self._flask.add_url_rule('/api/player_get_cover',
                             'player_get_cover',self.get_cover)
    self._flask.add_url_rule('/api/update_state','update_state',
                             self.update_state,methods=['POST'])
    self._flask.add_url_rule('/api/<path:api>','api',self.process_api)

  # --- return absolute path of web-files   ----------------------------------

  def _get_path(self,*path):
    """ absolute path of web-file """
    return os.path.join(self._web_root,*path)

  # --- static routes   ------------------------------------------------------

  def css_pages(self,filepath):
    return send_from_directory(self._get_path('css'),filepath)

  def webfonts(self,filepath):
    return send_from_directory(self._get_path('webfonts'),filepath)
  
  def images(self,filepath):
    return send_from_directory(self._get_path('images'),filepath)
  
  def js_pages(self,filepath):
    return send_from_directory(self._get_path('js'),filepath)
  
  # --- main page   ----------------------------------------------------------

  def main_page(self):
    return render_template("index.html")

  # --- process API-call   -------------------------------------------------

  def process_api(self,api):
    """ process Bhojpur Speech API """

    if api.startswith("_"):
      # internal API, illegal request!
      self.msg("illegal api-call: %s" % api)
      msg = '"illegal request /api/%s"' % api
      response = make_response(('{"msg": ' + msg +'}',400))
      response.content_type = 'application/json'
      return response
    else:
      self.msg("processing api-call: %s" % api)
      try:
        response = self._api._exec(api,**request.args)
        return json.dumps(response)
      except NotImplementedError as err:
        self.msg("illegal request: /api/%s" % api)
        msg = '"/api/%s not implemented"' % api
        response = make_response(('{"msg": ' + msg +'}',400))
        response.content_type = 'application/json'
        return response
      except Exception as ex:
        self.msg("exception while calling: /api/%s" % api)
        traceback.print_exc()
        msg = '"internal server error"'
        response = make_response(('{"msg": ' + msg +'}',500))
        response.content_type = 'application/json'
        return response

  # --- publish state   ----------------------------------------------------

  def update_state(self):
    """ update state and redistribute """

    try:
      # only a subset of the state is controlled by the client, so filter
      # for valid values
      state = request.get_json(force=True)
      for k in list(state.keys()):
        if k not in ['webgui','mode']:
          del state[k]
      self._api.update_state(state=state)
      return ""
    except:
      self.msg("Bhojpur Speech: exception while calling: /api/update_state")
      traceback.print_exc()
      msg = '"internal server error"'
      response = make_response(('{"msg": ' + msg +'}',500))
      response.content_type = 'application/json'
      return response

  # --- return cover   -----------------------------------------------------

  def get_cover(self,dir="ignored"):
    """ return cover if available """

    cover = self._api._player_get_cover_file()
    if cover:
      return send_file(cover)
    else:
      return send_from_directory(self._get_path('images'),'default.png')

  # --- stream SSE (server sent events)   ----------------------------------

  def get_events(self):
    """ stream SSE """

    try:
      id = uuid.uuid4().hex
      ev_queue = self._api._add_consumer(id)   # TODO: use session-id

      def event_stream():
        while True:
          ev = ev_queue.get()
          ev_queue.task_done()
          if ev:
            sse = "data: %s\n\n" % json.dumps(ev)
            #self.msg("WebServer: serving event '%s'" % sse)
            yield sse
          else:
            break

      return Response(event_stream(), mimetype='text/event-stream')
    except:
      traceback.print_exc()

  # --- stop web-server   --------------------------------------------------

  def stop(self):
    """ stop the Bhojpur Speech web-server """

    self.msg("WebServer: process stop-request")
    self._server.shutdown()

  # --- service-loop   -----------------------------------------------------

  def run(self):
    """ start and run the Bhojpur Speech webserver """

    class QuietHandler(WSGIRequestHandler):
      def log_request(*args, **kw):
        pass

    if self.debug:
      self._server = make_server(self._host,self._port,self._flask,threaded=True)
    else:
      self._server = make_server(self._host,self._port,self._flask,
                                 request_handler=QuietHandler,threaded=True)
    ctx = self._flask.app_context()
    ctx.push()

    self.msg("WebServer: starting the Bhojpur Speech web-server in debug-mode")
    self.msg("WebServer: listening on port %s" % self._port)
    self.msg("WebServer: using web-root: %s" % self._web_root)
    self._server.serve_forever()
    self.msg("WebServer: Bhojpur Speech finished")