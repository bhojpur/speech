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

#ifndef VOSK_KALDI_RECOGNIZER_H
#define VOSK_KALDI_RECOGNIZER_H

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

#include "model.h"
#include "spk_model.h"

using namespace kaldi;

enum RecognizerState {
    RECOGNIZER_INITIALIZED,
    RECOGNIZER_RUNNING,
    RECOGNIZER_ENDPOINT,
    RECOGNIZER_FINALIZED
};

class Recognizer {
    public:
        Recognizer(Model *model, float sample_frequency);
        Recognizer(Model *model, float sample_frequency, SpkModel *spk_model);
        Recognizer(Model *model, float sample_frequency, char const *grammar);
        ~Recognizer();
        void SetMaxAlternatives(int max_alternatives);
        void SetSpkModel(SpkModel *spk_model);
        void SetWords(bool words);
        void SetPartialWords(bool partial_words);
        void SetNLSML(bool nlsml);
        bool AcceptWaveform(const char *data, int len);
        bool AcceptWaveform(const short *sdata, int len);
        bool AcceptWaveform(const float *fdata, int len);
        const char* Result();
        const char* FinalResult();
        const char* PartialResult();
        void Reset();

    private:
        void InitState();
        void InitRescoring();
        void CleanUp();
        void UpdateSilenceWeights();
        bool AcceptWaveform(Vector<BaseFloat> &wdata);
        bool GetSpkVector(Vector<BaseFloat> &out_xvector, int *frames);
        const char *GetResult();
        const char *StoreEmptyReturn();
        const char *StoreReturn(const string &res);
        const char *MbrResult(CompactLattice &clat);
        const char *NbestResult(CompactLattice &clat);
        const char *NlsmlResult(CompactLattice &clat);

        Model *model_ = nullptr;
        SingleUtteranceNnet3IncrementalDecoder *decoder_ = nullptr;
        fst::LookaheadFst<fst::StdArc, int32> *decode_fst_ = nullptr;
        fst::StdVectorFst *g_fst_ = nullptr; // dynamically constructed grammar
        OnlineNnet2FeaturePipeline *feature_pipeline_ = nullptr;
        OnlineSilenceWeighting *silence_weighting_ = nullptr;

        // Speaker identification
        SpkModel *spk_model_ = nullptr;
        OnlineBaseFeature *spk_feature_ = nullptr;

        // Rescoring
        fst::ArcMapFst<fst::StdArc, LatticeArc, fst::StdToLatticeMapper<BaseFloat> > *lm_to_subtract_ = nullptr;
        kaldi::ConstArpaLmDeterministicFst *carpa_to_add_ = nullptr;
        fst::ScaleDeterministicOnDemandFst *carpa_to_add_scale_ = nullptr;
        // RNNLM rescoring
        kaldi::rnnlm::KaldiRnnlmDeterministicFst* rnnlm_to_add_ = nullptr;
        fst::DeterministicOnDemandFst<fst::StdArc> *rnnlm_to_add_scale_ = nullptr;
        kaldi::rnnlm::RnnlmComputeStateInfo *rnnlm_info_ = nullptr;


        // Other
        int max_alternatives_ = 0; // Disable alternatives by default
        bool words_ = false;
        bool partial_words_ = false;
        bool nlsml_ = false;

        float sample_frequency_;
        int32 frame_offset_;

        int64 samples_processed_;
        int64 samples_round_start_;

        RecognizerState state_;
        string last_result_;
};

#endif /* VOSK_KALDI_RECOGNIZER_H */