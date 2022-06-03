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

#ifdef __cplusplus
extern "C" {
#endif
    typedef struct TokenMetadata {
        const char* text;
        const unsigned int timestep;
        const float start_time;
    } TokenMetadata;

    typedef struct CandidateTranscript {
        const TokenMetadata* const tokens;
        const unsigned int num_tokens;
        const double confidence;
    } CandidateTranscript;

    typedef struct Metadata {
        const CandidateTranscript* const transcripts;
        const unsigned int num_transcripts;
    } Metadata;

    typedef void* ModelWrapper;
    ModelWrapper* New(const char* aModelPath, int* errorOut);
    void Model_Close(ModelWrapper* w);
    unsigned int Model_BeamWidth(ModelWrapper* w);
    int Model_SetBeamWidth(ModelWrapper* w, unsigned int aBeamWidth);
    int Model_SampleRate(ModelWrapper* w);
    int Model_EnableExternalScorer(ModelWrapper* w, const char* aScorerPath);
    int Model_DisableExternalScorer(ModelWrapper* w);
    int Model_SetScorerAlphaBeta(ModelWrapper* w, float aAlpha, float aBeta);
    char* Model_STT(ModelWrapper* w, const short* aBuffer, unsigned int aBufferSize);
    Metadata* Model_STTWithMetadata(ModelWrapper* w, const short* aBuffer, unsigned int aBufferSize, unsigned int aNumResults);

    typedef void* StreamWrapper;
    StreamWrapper* Model_NewStream(ModelWrapper* w, int* errorOut);
    void Stream_Discard(StreamWrapper* sw);
    void Stream_FeedAudioContent(StreamWrapper* sw, const short* aBuffer, unsigned int aBufferSize);
    char* Stream_IntermediateDecode(StreamWrapper* sw);
    Metadata* Stream_IntermediateDecodeWithMetadata(StreamWrapper* sw, unsigned int aNumResults);
    char* Stream_Finish(StreamWrapper* sw);
    Metadata* Stream_FinishWithMetadata(StreamWrapper* sw, unsigned int aNumResults);

    const CandidateTranscript* Metadata_Transcripts(Metadata* m);
    unsigned int Metadata_NumTranscripts(Metadata* m);
    void Metadata_Close(Metadata* m);

    const TokenMetadata* CandidateTranscript_Tokens(CandidateTranscript* ct);
    unsigned int CandidateTranscript_NumTokens(CandidateTranscript* ct);
    double CandidateTranscript_Confidence(CandidateTranscript* ct);

    const char* TokenMetadata_Text(TokenMetadata* tm);
    unsigned int TokenMetadata_Timestep(TokenMetadata* tm);
    float TokenMetadata_StartTime(TokenMetadata* tm);

    void FreeString(char* s);
    char* Version();
    char* ErrorCodeToErrorMessage(int aErrorCode);

#ifdef __cplusplus
}
#endif