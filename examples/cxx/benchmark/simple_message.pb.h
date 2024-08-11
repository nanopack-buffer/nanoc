// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: simple_message.proto

#ifndef GOOGLE_PROTOBUF_INCLUDED_simple_5fmessage_2eproto
#define GOOGLE_PROTOBUF_INCLUDED_simple_5fmessage_2eproto

#include <limits>
#include <string>

#include <google/protobuf/port_def.inc>
#if PROTOBUF_VERSION < 3021000
#error This file was generated by a newer version of protoc which is
#error incompatible with your Protocol Buffer headers. Please update
#error your headers.
#endif
#if 3021012 < PROTOBUF_MIN_PROTOC_VERSION
#error This file was generated by an older version of protoc which is
#error incompatible with your Protocol Buffer headers. Please
#error regenerate this file with a newer version of protoc.
#endif

#include <google/protobuf/port_undef.inc>
#include <google/protobuf/io/coded_stream.h>
#include <google/protobuf/arena.h>
#include <google/protobuf/arenastring.h>
#include <google/protobuf/generated_message_util.h>
#include <google/protobuf/metadata_lite.h>
#include <google/protobuf/generated_message_reflection.h>
#include <google/protobuf/message.h>
#include <google/protobuf/repeated_field.h>  // IWYU pragma: export
#include <google/protobuf/extension_set.h>  // IWYU pragma: export
#include <google/protobuf/map.h>  // IWYU pragma: export
#include <google/protobuf/map_entry.h>
#include <google/protobuf/map_field_inl.h>
#include <google/protobuf/unknown_field_set.h>
// @@protoc_insertion_point(includes)
#include <google/protobuf/port_def.inc>
#define PROTOBUF_INTERNAL_EXPORT_simple_5fmessage_2eproto
PROTOBUF_NAMESPACE_OPEN
namespace internal {
class AnyMetadata;
}  // namespace internal
PROTOBUF_NAMESPACE_CLOSE

// Internal implementation detail -- do not use these members.
struct TableStruct_simple_5fmessage_2eproto {
  static const uint32_t offsets[];
};
extern const ::PROTOBUF_NAMESPACE_ID::internal::DescriptorTable descriptor_table_simple_5fmessage_2eproto;
class PbSimpleMessage;
struct PbSimpleMessageDefaultTypeInternal;
extern PbSimpleMessageDefaultTypeInternal _PbSimpleMessage_default_instance_;
class PbSimpleMessage_MapFieldEntry_DoNotUse;
struct PbSimpleMessage_MapFieldEntry_DoNotUseDefaultTypeInternal;
extern PbSimpleMessage_MapFieldEntry_DoNotUseDefaultTypeInternal _PbSimpleMessage_MapFieldEntry_DoNotUse_default_instance_;
PROTOBUF_NAMESPACE_OPEN
template<> ::PbSimpleMessage* Arena::CreateMaybeMessage<::PbSimpleMessage>(Arena*);
template<> ::PbSimpleMessage_MapFieldEntry_DoNotUse* Arena::CreateMaybeMessage<::PbSimpleMessage_MapFieldEntry_DoNotUse>(Arena*);
PROTOBUF_NAMESPACE_CLOSE

// ===================================================================

class PbSimpleMessage_MapFieldEntry_DoNotUse : public ::PROTOBUF_NAMESPACE_ID::internal::MapEntry<PbSimpleMessage_MapFieldEntry_DoNotUse, 
    std::string, bool,
    ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::TYPE_STRING,
    ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::TYPE_BOOL> {
public:
  typedef ::PROTOBUF_NAMESPACE_ID::internal::MapEntry<PbSimpleMessage_MapFieldEntry_DoNotUse, 
    std::string, bool,
    ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::TYPE_STRING,
    ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::TYPE_BOOL> SuperType;
  PbSimpleMessage_MapFieldEntry_DoNotUse();
  explicit PROTOBUF_CONSTEXPR PbSimpleMessage_MapFieldEntry_DoNotUse(
      ::PROTOBUF_NAMESPACE_ID::internal::ConstantInitialized);
  explicit PbSimpleMessage_MapFieldEntry_DoNotUse(::PROTOBUF_NAMESPACE_ID::Arena* arena);
  void MergeFrom(const PbSimpleMessage_MapFieldEntry_DoNotUse& other);
  static const PbSimpleMessage_MapFieldEntry_DoNotUse* internal_default_instance() { return reinterpret_cast<const PbSimpleMessage_MapFieldEntry_DoNotUse*>(&_PbSimpleMessage_MapFieldEntry_DoNotUse_default_instance_); }
  static bool ValidateKey(std::string* s) {
    return ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::VerifyUtf8String(s->data(), static_cast<int>(s->size()), ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::PARSE, "PbSimpleMessage.MapFieldEntry.key");
 }
  static bool ValidateValue(void*) { return true; }
  using ::PROTOBUF_NAMESPACE_ID::Message::MergeFrom;
  ::PROTOBUF_NAMESPACE_ID::Metadata GetMetadata() const final;
  friend struct ::TableStruct_simple_5fmessage_2eproto;
};

// -------------------------------------------------------------------

class PbSimpleMessage final :
    public ::PROTOBUF_NAMESPACE_ID::Message /* @@protoc_insertion_point(class_definition:PbSimpleMessage) */ {
 public:
  inline PbSimpleMessage() : PbSimpleMessage(nullptr) {}
  ~PbSimpleMessage() override;
  explicit PROTOBUF_CONSTEXPR PbSimpleMessage(::PROTOBUF_NAMESPACE_ID::internal::ConstantInitialized);

  PbSimpleMessage(const PbSimpleMessage& from);
  PbSimpleMessage(PbSimpleMessage&& from) noexcept
    : PbSimpleMessage() {
    *this = ::std::move(from);
  }

  inline PbSimpleMessage& operator=(const PbSimpleMessage& from) {
    CopyFrom(from);
    return *this;
  }
  inline PbSimpleMessage& operator=(PbSimpleMessage&& from) noexcept {
    if (this == &from) return *this;
    if (GetOwningArena() == from.GetOwningArena()
  #ifdef PROTOBUF_FORCE_COPY_IN_MOVE
        && GetOwningArena() != nullptr
  #endif  // !PROTOBUF_FORCE_COPY_IN_MOVE
    ) {
      InternalSwap(&from);
    } else {
      CopyFrom(from);
    }
    return *this;
  }

  static const ::PROTOBUF_NAMESPACE_ID::Descriptor* descriptor() {
    return GetDescriptor();
  }
  static const ::PROTOBUF_NAMESPACE_ID::Descriptor* GetDescriptor() {
    return default_instance().GetMetadata().descriptor;
  }
  static const ::PROTOBUF_NAMESPACE_ID::Reflection* GetReflection() {
    return default_instance().GetMetadata().reflection;
  }
  static const PbSimpleMessage& default_instance() {
    return *internal_default_instance();
  }
  static inline const PbSimpleMessage* internal_default_instance() {
    return reinterpret_cast<const PbSimpleMessage*>(
               &_PbSimpleMessage_default_instance_);
  }
  static constexpr int kIndexInFileMessages =
    1;

  friend void swap(PbSimpleMessage& a, PbSimpleMessage& b) {
    a.Swap(&b);
  }
  inline void Swap(PbSimpleMessage* other) {
    if (other == this) return;
  #ifdef PROTOBUF_FORCE_COPY_IN_SWAP
    if (GetOwningArena() != nullptr &&
        GetOwningArena() == other->GetOwningArena()) {
   #else  // PROTOBUF_FORCE_COPY_IN_SWAP
    if (GetOwningArena() == other->GetOwningArena()) {
  #endif  // !PROTOBUF_FORCE_COPY_IN_SWAP
      InternalSwap(other);
    } else {
      ::PROTOBUF_NAMESPACE_ID::internal::GenericSwap(this, other);
    }
  }
  void UnsafeArenaSwap(PbSimpleMessage* other) {
    if (other == this) return;
    GOOGLE_DCHECK(GetOwningArena() == other->GetOwningArena());
    InternalSwap(other);
  }

  // implements Message ----------------------------------------------

  PbSimpleMessage* New(::PROTOBUF_NAMESPACE_ID::Arena* arena = nullptr) const final {
    return CreateMaybeMessage<PbSimpleMessage>(arena);
  }
  using ::PROTOBUF_NAMESPACE_ID::Message::CopyFrom;
  void CopyFrom(const PbSimpleMessage& from);
  using ::PROTOBUF_NAMESPACE_ID::Message::MergeFrom;
  void MergeFrom( const PbSimpleMessage& from) {
    PbSimpleMessage::MergeImpl(*this, from);
  }
  private:
  static void MergeImpl(::PROTOBUF_NAMESPACE_ID::Message& to_msg, const ::PROTOBUF_NAMESPACE_ID::Message& from_msg);
  public:
  PROTOBUF_ATTRIBUTE_REINITIALIZES void Clear() final;
  bool IsInitialized() const final;

  size_t ByteSizeLong() const final;
  const char* _InternalParse(const char* ptr, ::PROTOBUF_NAMESPACE_ID::internal::ParseContext* ctx) final;
  uint8_t* _InternalSerialize(
      uint8_t* target, ::PROTOBUF_NAMESPACE_ID::io::EpsCopyOutputStream* stream) const final;
  int GetCachedSize() const final { return _impl_._cached_size_.Get(); }

  private:
  void SharedCtor(::PROTOBUF_NAMESPACE_ID::Arena* arena, bool is_message_owned);
  void SharedDtor();
  void SetCachedSize(int size) const final;
  void InternalSwap(PbSimpleMessage* other);

  private:
  friend class ::PROTOBUF_NAMESPACE_ID::internal::AnyMetadata;
  static ::PROTOBUF_NAMESPACE_ID::StringPiece FullMessageName() {
    return "PbSimpleMessage";
  }
  protected:
  explicit PbSimpleMessage(::PROTOBUF_NAMESPACE_ID::Arena* arena,
                       bool is_message_owned = false);
  private:
  static void ArenaDtor(void* object);
  public:

  static const ClassData _class_data_;
  const ::PROTOBUF_NAMESPACE_ID::Message::ClassData*GetClassData() const final;

  ::PROTOBUF_NAMESPACE_ID::Metadata GetMetadata() const final;

  // nested types ----------------------------------------------------


  // accessors -------------------------------------------------------

  enum : int {
    kArrayFieldFieldNumber = 5,
    kMapFieldFieldNumber = 6,
    kStringFieldFieldNumber = 1,
    kOptionalFieldFieldNumber = 4,
    kDoubleFieldFieldNumber = 3,
    kIntFieldFieldNumber = 2,
  };
  // repeated int32 array_field = 5;
  int array_field_size() const;
  private:
  int _internal_array_field_size() const;
  public:
  void clear_array_field();
  private:
  int32_t _internal_array_field(int index) const;
  const ::PROTOBUF_NAMESPACE_ID::RepeatedField< int32_t >&
      _internal_array_field() const;
  void _internal_add_array_field(int32_t value);
  ::PROTOBUF_NAMESPACE_ID::RepeatedField< int32_t >*
      _internal_mutable_array_field();
  public:
  int32_t array_field(int index) const;
  void set_array_field(int index, int32_t value);
  void add_array_field(int32_t value);
  const ::PROTOBUF_NAMESPACE_ID::RepeatedField< int32_t >&
      array_field() const;
  ::PROTOBUF_NAMESPACE_ID::RepeatedField< int32_t >*
      mutable_array_field();

  // map<string, bool> map_field = 6;
  int map_field_size() const;
  private:
  int _internal_map_field_size() const;
  public:
  void clear_map_field();
  private:
  const ::PROTOBUF_NAMESPACE_ID::Map< std::string, bool >&
      _internal_map_field() const;
  ::PROTOBUF_NAMESPACE_ID::Map< std::string, bool >*
      _internal_mutable_map_field();
  public:
  const ::PROTOBUF_NAMESPACE_ID::Map< std::string, bool >&
      map_field() const;
  ::PROTOBUF_NAMESPACE_ID::Map< std::string, bool >*
      mutable_map_field();

  // string string_field = 1;
  void clear_string_field();
  const std::string& string_field() const;
  template <typename ArgT0 = const std::string&, typename... ArgT>
  void set_string_field(ArgT0&& arg0, ArgT... args);
  std::string* mutable_string_field();
  PROTOBUF_NODISCARD std::string* release_string_field();
  void set_allocated_string_field(std::string* string_field);
  private:
  const std::string& _internal_string_field() const;
  inline PROTOBUF_ALWAYS_INLINE void _internal_set_string_field(const std::string& value);
  std::string* _internal_mutable_string_field();
  public:

  // string optional_field = 4;
  void clear_optional_field();
  const std::string& optional_field() const;
  template <typename ArgT0 = const std::string&, typename... ArgT>
  void set_optional_field(ArgT0&& arg0, ArgT... args);
  std::string* mutable_optional_field();
  PROTOBUF_NODISCARD std::string* release_optional_field();
  void set_allocated_optional_field(std::string* optional_field);
  private:
  const std::string& _internal_optional_field() const;
  inline PROTOBUF_ALWAYS_INLINE void _internal_set_optional_field(const std::string& value);
  std::string* _internal_mutable_optional_field();
  public:

  // double double_field = 3;
  void clear_double_field();
  double double_field() const;
  void set_double_field(double value);
  private:
  double _internal_double_field() const;
  void _internal_set_double_field(double value);
  public:

  // int32 int_field = 2;
  void clear_int_field();
  int32_t int_field() const;
  void set_int_field(int32_t value);
  private:
  int32_t _internal_int_field() const;
  void _internal_set_int_field(int32_t value);
  public:

  // @@protoc_insertion_point(class_scope:PbSimpleMessage)
 private:
  class _Internal;

  template <typename T> friend class ::PROTOBUF_NAMESPACE_ID::Arena::InternalHelper;
  typedef void InternalArenaConstructable_;
  typedef void DestructorSkippable_;
  struct Impl_ {
    ::PROTOBUF_NAMESPACE_ID::RepeatedField< int32_t > array_field_;
    mutable std::atomic<int> _array_field_cached_byte_size_;
    ::PROTOBUF_NAMESPACE_ID::internal::MapField<
        PbSimpleMessage_MapFieldEntry_DoNotUse,
        std::string, bool,
        ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::TYPE_STRING,
        ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::TYPE_BOOL> map_field_;
    ::PROTOBUF_NAMESPACE_ID::internal::ArenaStringPtr string_field_;
    ::PROTOBUF_NAMESPACE_ID::internal::ArenaStringPtr optional_field_;
    double double_field_;
    int32_t int_field_;
    mutable ::PROTOBUF_NAMESPACE_ID::internal::CachedSize _cached_size_;
  };
  union { Impl_ _impl_; };
  friend struct ::TableStruct_simple_5fmessage_2eproto;
};
// ===================================================================


// ===================================================================

#ifdef __GNUC__
  #pragma GCC diagnostic push
  #pragma GCC diagnostic ignored "-Wstrict-aliasing"
#endif  // __GNUC__
// -------------------------------------------------------------------

// PbSimpleMessage

// string string_field = 1;
inline void PbSimpleMessage::clear_string_field() {
  _impl_.string_field_.ClearToEmpty();
}
inline const std::string& PbSimpleMessage::string_field() const {
  // @@protoc_insertion_point(field_get:PbSimpleMessage.string_field)
  return _internal_string_field();
}
template <typename ArgT0, typename... ArgT>
inline PROTOBUF_ALWAYS_INLINE
void PbSimpleMessage::set_string_field(ArgT0&& arg0, ArgT... args) {
 
 _impl_.string_field_.Set(static_cast<ArgT0 &&>(arg0), args..., GetArenaForAllocation());
  // @@protoc_insertion_point(field_set:PbSimpleMessage.string_field)
}
inline std::string* PbSimpleMessage::mutable_string_field() {
  std::string* _s = _internal_mutable_string_field();
  // @@protoc_insertion_point(field_mutable:PbSimpleMessage.string_field)
  return _s;
}
inline const std::string& PbSimpleMessage::_internal_string_field() const {
  return _impl_.string_field_.Get();
}
inline void PbSimpleMessage::_internal_set_string_field(const std::string& value) {
  
  _impl_.string_field_.Set(value, GetArenaForAllocation());
}
inline std::string* PbSimpleMessage::_internal_mutable_string_field() {
  
  return _impl_.string_field_.Mutable(GetArenaForAllocation());
}
inline std::string* PbSimpleMessage::release_string_field() {
  // @@protoc_insertion_point(field_release:PbSimpleMessage.string_field)
  return _impl_.string_field_.Release();
}
inline void PbSimpleMessage::set_allocated_string_field(std::string* string_field) {
  if (string_field != nullptr) {
    
  } else {
    
  }
  _impl_.string_field_.SetAllocated(string_field, GetArenaForAllocation());
#ifdef PROTOBUF_FORCE_COPY_DEFAULT_STRING
  if (_impl_.string_field_.IsDefault()) {
    _impl_.string_field_.Set("", GetArenaForAllocation());
  }
#endif // PROTOBUF_FORCE_COPY_DEFAULT_STRING
  // @@protoc_insertion_point(field_set_allocated:PbSimpleMessage.string_field)
}

// int32 int_field = 2;
inline void PbSimpleMessage::clear_int_field() {
  _impl_.int_field_ = 0;
}
inline int32_t PbSimpleMessage::_internal_int_field() const {
  return _impl_.int_field_;
}
inline int32_t PbSimpleMessage::int_field() const {
  // @@protoc_insertion_point(field_get:PbSimpleMessage.int_field)
  return _internal_int_field();
}
inline void PbSimpleMessage::_internal_set_int_field(int32_t value) {
  
  _impl_.int_field_ = value;
}
inline void PbSimpleMessage::set_int_field(int32_t value) {
  _internal_set_int_field(value);
  // @@protoc_insertion_point(field_set:PbSimpleMessage.int_field)
}

// double double_field = 3;
inline void PbSimpleMessage::clear_double_field() {
  _impl_.double_field_ = 0;
}
inline double PbSimpleMessage::_internal_double_field() const {
  return _impl_.double_field_;
}
inline double PbSimpleMessage::double_field() const {
  // @@protoc_insertion_point(field_get:PbSimpleMessage.double_field)
  return _internal_double_field();
}
inline void PbSimpleMessage::_internal_set_double_field(double value) {
  
  _impl_.double_field_ = value;
}
inline void PbSimpleMessage::set_double_field(double value) {
  _internal_set_double_field(value);
  // @@protoc_insertion_point(field_set:PbSimpleMessage.double_field)
}

// string optional_field = 4;
inline void PbSimpleMessage::clear_optional_field() {
  _impl_.optional_field_.ClearToEmpty();
}
inline const std::string& PbSimpleMessage::optional_field() const {
  // @@protoc_insertion_point(field_get:PbSimpleMessage.optional_field)
  return _internal_optional_field();
}
template <typename ArgT0, typename... ArgT>
inline PROTOBUF_ALWAYS_INLINE
void PbSimpleMessage::set_optional_field(ArgT0&& arg0, ArgT... args) {
 
 _impl_.optional_field_.Set(static_cast<ArgT0 &&>(arg0), args..., GetArenaForAllocation());
  // @@protoc_insertion_point(field_set:PbSimpleMessage.optional_field)
}
inline std::string* PbSimpleMessage::mutable_optional_field() {
  std::string* _s = _internal_mutable_optional_field();
  // @@protoc_insertion_point(field_mutable:PbSimpleMessage.optional_field)
  return _s;
}
inline const std::string& PbSimpleMessage::_internal_optional_field() const {
  return _impl_.optional_field_.Get();
}
inline void PbSimpleMessage::_internal_set_optional_field(const std::string& value) {
  
  _impl_.optional_field_.Set(value, GetArenaForAllocation());
}
inline std::string* PbSimpleMessage::_internal_mutable_optional_field() {
  
  return _impl_.optional_field_.Mutable(GetArenaForAllocation());
}
inline std::string* PbSimpleMessage::release_optional_field() {
  // @@protoc_insertion_point(field_release:PbSimpleMessage.optional_field)
  return _impl_.optional_field_.Release();
}
inline void PbSimpleMessage::set_allocated_optional_field(std::string* optional_field) {
  if (optional_field != nullptr) {
    
  } else {
    
  }
  _impl_.optional_field_.SetAllocated(optional_field, GetArenaForAllocation());
#ifdef PROTOBUF_FORCE_COPY_DEFAULT_STRING
  if (_impl_.optional_field_.IsDefault()) {
    _impl_.optional_field_.Set("", GetArenaForAllocation());
  }
#endif // PROTOBUF_FORCE_COPY_DEFAULT_STRING
  // @@protoc_insertion_point(field_set_allocated:PbSimpleMessage.optional_field)
}

// repeated int32 array_field = 5;
inline int PbSimpleMessage::_internal_array_field_size() const {
  return _impl_.array_field_.size();
}
inline int PbSimpleMessage::array_field_size() const {
  return _internal_array_field_size();
}
inline void PbSimpleMessage::clear_array_field() {
  _impl_.array_field_.Clear();
}
inline int32_t PbSimpleMessage::_internal_array_field(int index) const {
  return _impl_.array_field_.Get(index);
}
inline int32_t PbSimpleMessage::array_field(int index) const {
  // @@protoc_insertion_point(field_get:PbSimpleMessage.array_field)
  return _internal_array_field(index);
}
inline void PbSimpleMessage::set_array_field(int index, int32_t value) {
  _impl_.array_field_.Set(index, value);
  // @@protoc_insertion_point(field_set:PbSimpleMessage.array_field)
}
inline void PbSimpleMessage::_internal_add_array_field(int32_t value) {
  _impl_.array_field_.Add(value);
}
inline void PbSimpleMessage::add_array_field(int32_t value) {
  _internal_add_array_field(value);
  // @@protoc_insertion_point(field_add:PbSimpleMessage.array_field)
}
inline const ::PROTOBUF_NAMESPACE_ID::RepeatedField< int32_t >&
PbSimpleMessage::_internal_array_field() const {
  return _impl_.array_field_;
}
inline const ::PROTOBUF_NAMESPACE_ID::RepeatedField< int32_t >&
PbSimpleMessage::array_field() const {
  // @@protoc_insertion_point(field_list:PbSimpleMessage.array_field)
  return _internal_array_field();
}
inline ::PROTOBUF_NAMESPACE_ID::RepeatedField< int32_t >*
PbSimpleMessage::_internal_mutable_array_field() {
  return &_impl_.array_field_;
}
inline ::PROTOBUF_NAMESPACE_ID::RepeatedField< int32_t >*
PbSimpleMessage::mutable_array_field() {
  // @@protoc_insertion_point(field_mutable_list:PbSimpleMessage.array_field)
  return _internal_mutable_array_field();
}

// map<string, bool> map_field = 6;
inline int PbSimpleMessage::_internal_map_field_size() const {
  return _impl_.map_field_.size();
}
inline int PbSimpleMessage::map_field_size() const {
  return _internal_map_field_size();
}
inline void PbSimpleMessage::clear_map_field() {
  _impl_.map_field_.Clear();
}
inline const ::PROTOBUF_NAMESPACE_ID::Map< std::string, bool >&
PbSimpleMessage::_internal_map_field() const {
  return _impl_.map_field_.GetMap();
}
inline const ::PROTOBUF_NAMESPACE_ID::Map< std::string, bool >&
PbSimpleMessage::map_field() const {
  // @@protoc_insertion_point(field_map:PbSimpleMessage.map_field)
  return _internal_map_field();
}
inline ::PROTOBUF_NAMESPACE_ID::Map< std::string, bool >*
PbSimpleMessage::_internal_mutable_map_field() {
  return _impl_.map_field_.MutableMap();
}
inline ::PROTOBUF_NAMESPACE_ID::Map< std::string, bool >*
PbSimpleMessage::mutable_map_field() {
  // @@protoc_insertion_point(field_mutable_map:PbSimpleMessage.map_field)
  return _internal_mutable_map_field();
}

#ifdef __GNUC__
  #pragma GCC diagnostic pop
#endif  // __GNUC__
// -------------------------------------------------------------------


// @@protoc_insertion_point(namespace_scope)


// @@protoc_insertion_point(global_scope)

#include <google/protobuf/port_undef.inc>
#endif  // GOOGLE_PROTOBUF_INCLUDED_GOOGLE_PROTOBUF_INCLUDED_simple_5fmessage_2eproto
