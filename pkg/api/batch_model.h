// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

#ifndef VOSK_BATCH_MODEL_H
#define VOSK_BATCH_MODEL_H

#include "base/kaldi-common.h"
#include "util/common-utils.h"
#include "fstext/fstext-lib.h"
#include "fstext/fstext-utils.h"
#include "decoder/lattice-faster-decoder.h"
#include "feat/feature-mfcc.h"
#include "lat/kaldi-lattice.h"
#include "lat/word-align-lattice.h"
#include "lat/compose-lattice-pruned.h"
#include "nnet3/am-nnet-simple.h"
#include "nnet3/nnet-am-decodable-simple.h"
#include "nnet3/nnet-utils.h"

#include "cudadecoder/cuda-online-pipeline-dynamic-batcher.h"
#include "cudadecoder/batched-threaded-nnet3-cuda-online-pipeline.h"
#include "cudadecoder/batched-threaded-nnet3-cuda-pipeline2.h"
#include "cudadecoder/cuda-pipeline-common.h"

#include "model.h"

using namespace kaldi;
using namespace kaldi::cuda_decoder;

class BatchRecognizer;

class BatchModel {
    public:
        BatchModel();
        ~BatchModel();

        uint64_t GetID(BatchRecognizer *recognizer);
        void WaitForCompletion();

    private:
        friend class BatchRecognizer;

        kaldi::TransitionModel *trans_model_ = nullptr;
        kaldi::nnet3::AmNnetSimple *nnet_ = nullptr;
        const fst::SymbolTable *word_syms_ = nullptr;

        fst::Fst<fst::StdArc> *hclg_fst_ = nullptr;
        kaldi::WordBoundaryInfo *winfo_ = nullptr;

        BatchedThreadedNnet3CudaOnlinePipeline *cuda_pipeline_ = nullptr;
        CudaOnlinePipelineDynamicBatcher *dynamic_batcher_ = nullptr;

        int32 samples_per_chunk_;
        uint64_t last_id_;
};

#endif /* VOSK_BATCH_MODEL_H */