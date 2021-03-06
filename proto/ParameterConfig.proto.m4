/* Copyright (c) 2016 PaddlePaddle Authors. All Rights Reserve.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License. */
ifdef(`proto3', `syntax = "proto2";')

package paddle;

/**
 * Configuration structure for parameter
 */

enum ParameterInitStrategy {
  PARAMETER_INIT_NORMAL = 0;
  PARAMETER_INIT_UNIFORM = 1;
}

message ParameterUpdaterHookConfig {
  required string type = 1;
  optional string purning_mask_filename = 2;
}

message ParameterConfig {
  required string name = 1;
  required uint64 size = 2;
  optional real learning_rate = 3 [default = 1.0];
  optional real momentum = 4 [default = 0.0];
  optional real initial_mean = 5 [default = 0.0];
  optional real initial_std = 6 [default = 0.01];
  // use L2-regularization if decay_rate set and decay_rate_l1 not set
  optional real decay_rate = 7 [default = 0.0];
  // use L1-regularization if decay_rate_l1 set
  optional real decay_rate_l1 = 8 [default = 0.0];
  // dims of Parameter, e.g. dims[0] as height, dims[1] as width..
  repeated uint64 dims = 9;
  // the gpu device which the parameter in.
  // Only used by ParallelNeuralNetork. Ignored otherwise.
  optional int32 device = 10 [default = -1];
  // how to init the parameter: 0 -> normal, 1 -> uniform
  // 0: treat initial_mean as mean, intial_std as standard deviation
  // 1: range is (initial_mean - initial_std) to (initial_mean + initial_std)
  optional int32 initial_strategy = 11 [default = 0];
  // define the variance when init the parameter, by height of the Matrix
  optional bool initial_smart = 12 [default = false];
  // apply regularization every # batches
  optional int32 num_batches_regularization = 13 [default = 1];
  // if is_sparse is true, para is sparse, else para is dense
  optional bool is_sparse = 14[default = false];
  // if para is sparse, format should be "csc" or "csr", empty means is not sparse
  optional string format = 15 [default = ""];
  // sparse remote update or not
  optional bool sparse_remote_update = 16 [default = false];
  // gradient clipping threshold, no clipping by default
  optional real gradient_clipping_threshold = 17 [default = 0.0];
  // static parameters are fixed when training
  optional bool is_static = 18 [default = false];
  // para_id should NOT be set by config_parser. It is for
  // internal use.
  optional uint64 para_id = 19;

  repeated ParameterUpdaterHookConfig update_hooks = 20;
  // setup load mat -> csr
  optional bool need_compact = 21 [default = false];
  // whether to do sparse update for this parameter
  optional bool sparse_update = 22 [default = false];

  // whether this parameter is shared or not.
  optional bool is_shared = 23 [default = false];
  // parameter block size
  optional uint64 parameter_block_size = 24 [default = 0];
}
