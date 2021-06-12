// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides

part of 'recipes.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

T _$identity<T>(T value) => value;

final _privateConstructorUsedError = UnsupportedError(
    'It seems like you constructed your class using `MyClass._()`. This constructor is only meant to be used by freezed and you are not supposed to need it nor use it.\nPlease check the documentation here for more informations: https://github.com/rrousselGit/freezed#custom-getters-and-methods');

/// @nodoc
class _$RecipesStateTearOff {
  const _$RecipesStateTearOff();

  _RecipesStateInitial initial() {
    return const _RecipesStateInitial();
  }

  _RecipesStateLoading loading() {
    return const _RecipesStateLoading();
  }

  _RecipesStateData data({required RecipesModel recipes}) {
    return _RecipesStateData(
      recipes: recipes,
    );
  }

  _RecipesStateError error([String? error]) {
    return _RecipesStateError(
      error,
    );
  }
}

/// @nodoc
const $RecipesState = _$RecipesStateTearOff();

/// @nodoc
mixin _$RecipesState {
  @optionalTypeArgs
  TResult when<TResult extends Object?>({
    required TResult Function() initial,
    required TResult Function() loading,
    required TResult Function(RecipesModel recipes) data,
    required TResult Function(String? error) error,
  }) =>
      throw _privateConstructorUsedError;
  @optionalTypeArgs
  TResult maybeWhen<TResult extends Object?>({
    TResult Function()? initial,
    TResult Function()? loading,
    TResult Function(RecipesModel recipes)? data,
    TResult Function(String? error)? error,
    required TResult orElse(),
  }) =>
      throw _privateConstructorUsedError;
  @optionalTypeArgs
  TResult map<TResult extends Object?>({
    required TResult Function(_RecipesStateInitial value) initial,
    required TResult Function(_RecipesStateLoading value) loading,
    required TResult Function(_RecipesStateData value) data,
    required TResult Function(_RecipesStateError value) error,
  }) =>
      throw _privateConstructorUsedError;
  @optionalTypeArgs
  TResult maybeMap<TResult extends Object?>({
    TResult Function(_RecipesStateInitial value)? initial,
    TResult Function(_RecipesStateLoading value)? loading,
    TResult Function(_RecipesStateData value)? data,
    TResult Function(_RecipesStateError value)? error,
    required TResult orElse(),
  }) =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $RecipesStateCopyWith<$Res> {
  factory $RecipesStateCopyWith(
          RecipesState value, $Res Function(RecipesState) then) =
      _$RecipesStateCopyWithImpl<$Res>;
}

/// @nodoc
class _$RecipesStateCopyWithImpl<$Res> implements $RecipesStateCopyWith<$Res> {
  _$RecipesStateCopyWithImpl(this._value, this._then);

  final RecipesState _value;
  // ignore: unused_field
  final $Res Function(RecipesState) _then;
}

/// @nodoc
abstract class _$RecipesStateInitialCopyWith<$Res> {
  factory _$RecipesStateInitialCopyWith(_RecipesStateInitial value,
          $Res Function(_RecipesStateInitial) then) =
      __$RecipesStateInitialCopyWithImpl<$Res>;
}

/// @nodoc
class __$RecipesStateInitialCopyWithImpl<$Res>
    extends _$RecipesStateCopyWithImpl<$Res>
    implements _$RecipesStateInitialCopyWith<$Res> {
  __$RecipesStateInitialCopyWithImpl(
      _RecipesStateInitial _value, $Res Function(_RecipesStateInitial) _then)
      : super(_value, (v) => _then(v as _RecipesStateInitial));

  @override
  _RecipesStateInitial get _value => super._value as _RecipesStateInitial;
}

/// @nodoc

class _$_RecipesStateInitial implements _RecipesStateInitial {
  const _$_RecipesStateInitial();

  @override
  String toString() {
    return 'RecipesState.initial()';
  }

  @override
  bool operator ==(dynamic other) {
    return identical(this, other) || (other is _RecipesStateInitial);
  }

  @override
  int get hashCode => runtimeType.hashCode;

  @override
  @optionalTypeArgs
  TResult when<TResult extends Object?>({
    required TResult Function() initial,
    required TResult Function() loading,
    required TResult Function(RecipesModel recipes) data,
    required TResult Function(String? error) error,
  }) {
    return initial();
  }

  @override
  @optionalTypeArgs
  TResult maybeWhen<TResult extends Object?>({
    TResult Function()? initial,
    TResult Function()? loading,
    TResult Function(RecipesModel recipes)? data,
    TResult Function(String? error)? error,
    required TResult orElse(),
  }) {
    if (initial != null) {
      return initial();
    }
    return orElse();
  }

  @override
  @optionalTypeArgs
  TResult map<TResult extends Object?>({
    required TResult Function(_RecipesStateInitial value) initial,
    required TResult Function(_RecipesStateLoading value) loading,
    required TResult Function(_RecipesStateData value) data,
    required TResult Function(_RecipesStateError value) error,
  }) {
    return initial(this);
  }

  @override
  @optionalTypeArgs
  TResult maybeMap<TResult extends Object?>({
    TResult Function(_RecipesStateInitial value)? initial,
    TResult Function(_RecipesStateLoading value)? loading,
    TResult Function(_RecipesStateData value)? data,
    TResult Function(_RecipesStateError value)? error,
    required TResult orElse(),
  }) {
    if (initial != null) {
      return initial(this);
    }
    return orElse();
  }
}

abstract class _RecipesStateInitial implements RecipesState {
  const factory _RecipesStateInitial() = _$_RecipesStateInitial;
}

/// @nodoc
abstract class _$RecipesStateLoadingCopyWith<$Res> {
  factory _$RecipesStateLoadingCopyWith(_RecipesStateLoading value,
          $Res Function(_RecipesStateLoading) then) =
      __$RecipesStateLoadingCopyWithImpl<$Res>;
}

/// @nodoc
class __$RecipesStateLoadingCopyWithImpl<$Res>
    extends _$RecipesStateCopyWithImpl<$Res>
    implements _$RecipesStateLoadingCopyWith<$Res> {
  __$RecipesStateLoadingCopyWithImpl(
      _RecipesStateLoading _value, $Res Function(_RecipesStateLoading) _then)
      : super(_value, (v) => _then(v as _RecipesStateLoading));

  @override
  _RecipesStateLoading get _value => super._value as _RecipesStateLoading;
}

/// @nodoc

class _$_RecipesStateLoading implements _RecipesStateLoading {
  const _$_RecipesStateLoading();

  @override
  String toString() {
    return 'RecipesState.loading()';
  }

  @override
  bool operator ==(dynamic other) {
    return identical(this, other) || (other is _RecipesStateLoading);
  }

  @override
  int get hashCode => runtimeType.hashCode;

  @override
  @optionalTypeArgs
  TResult when<TResult extends Object?>({
    required TResult Function() initial,
    required TResult Function() loading,
    required TResult Function(RecipesModel recipes) data,
    required TResult Function(String? error) error,
  }) {
    return loading();
  }

  @override
  @optionalTypeArgs
  TResult maybeWhen<TResult extends Object?>({
    TResult Function()? initial,
    TResult Function()? loading,
    TResult Function(RecipesModel recipes)? data,
    TResult Function(String? error)? error,
    required TResult orElse(),
  }) {
    if (loading != null) {
      return loading();
    }
    return orElse();
  }

  @override
  @optionalTypeArgs
  TResult map<TResult extends Object?>({
    required TResult Function(_RecipesStateInitial value) initial,
    required TResult Function(_RecipesStateLoading value) loading,
    required TResult Function(_RecipesStateData value) data,
    required TResult Function(_RecipesStateError value) error,
  }) {
    return loading(this);
  }

  @override
  @optionalTypeArgs
  TResult maybeMap<TResult extends Object?>({
    TResult Function(_RecipesStateInitial value)? initial,
    TResult Function(_RecipesStateLoading value)? loading,
    TResult Function(_RecipesStateData value)? data,
    TResult Function(_RecipesStateError value)? error,
    required TResult orElse(),
  }) {
    if (loading != null) {
      return loading(this);
    }
    return orElse();
  }
}

abstract class _RecipesStateLoading implements RecipesState {
  const factory _RecipesStateLoading() = _$_RecipesStateLoading;
}

/// @nodoc
abstract class _$RecipesStateDataCopyWith<$Res> {
  factory _$RecipesStateDataCopyWith(
          _RecipesStateData value, $Res Function(_RecipesStateData) then) =
      __$RecipesStateDataCopyWithImpl<$Res>;
  $Res call({RecipesModel recipes});
}

/// @nodoc
class __$RecipesStateDataCopyWithImpl<$Res>
    extends _$RecipesStateCopyWithImpl<$Res>
    implements _$RecipesStateDataCopyWith<$Res> {
  __$RecipesStateDataCopyWithImpl(
      _RecipesStateData _value, $Res Function(_RecipesStateData) _then)
      : super(_value, (v) => _then(v as _RecipesStateData));

  @override
  _RecipesStateData get _value => super._value as _RecipesStateData;

  @override
  $Res call({
    Object? recipes = freezed,
  }) {
    return _then(_RecipesStateData(
      recipes: recipes == freezed
          ? _value.recipes
          : recipes // ignore: cast_nullable_to_non_nullable
              as RecipesModel,
    ));
  }
}

/// @nodoc

class _$_RecipesStateData implements _RecipesStateData {
  const _$_RecipesStateData({required this.recipes});

  @override
  final RecipesModel recipes;

  @override
  String toString() {
    return 'RecipesState.data(recipes: $recipes)';
  }

  @override
  bool operator ==(dynamic other) {
    return identical(this, other) ||
        (other is _RecipesStateData &&
            (identical(other.recipes, recipes) ||
                const DeepCollectionEquality().equals(other.recipes, recipes)));
  }

  @override
  int get hashCode =>
      runtimeType.hashCode ^ const DeepCollectionEquality().hash(recipes);

  @JsonKey(ignore: true)
  @override
  _$RecipesStateDataCopyWith<_RecipesStateData> get copyWith =>
      __$RecipesStateDataCopyWithImpl<_RecipesStateData>(this, _$identity);

  @override
  @optionalTypeArgs
  TResult when<TResult extends Object?>({
    required TResult Function() initial,
    required TResult Function() loading,
    required TResult Function(RecipesModel recipes) data,
    required TResult Function(String? error) error,
  }) {
    return data(recipes);
  }

  @override
  @optionalTypeArgs
  TResult maybeWhen<TResult extends Object?>({
    TResult Function()? initial,
    TResult Function()? loading,
    TResult Function(RecipesModel recipes)? data,
    TResult Function(String? error)? error,
    required TResult orElse(),
  }) {
    if (data != null) {
      return data(recipes);
    }
    return orElse();
  }

  @override
  @optionalTypeArgs
  TResult map<TResult extends Object?>({
    required TResult Function(_RecipesStateInitial value) initial,
    required TResult Function(_RecipesStateLoading value) loading,
    required TResult Function(_RecipesStateData value) data,
    required TResult Function(_RecipesStateError value) error,
  }) {
    return data(this);
  }

  @override
  @optionalTypeArgs
  TResult maybeMap<TResult extends Object?>({
    TResult Function(_RecipesStateInitial value)? initial,
    TResult Function(_RecipesStateLoading value)? loading,
    TResult Function(_RecipesStateData value)? data,
    TResult Function(_RecipesStateError value)? error,
    required TResult orElse(),
  }) {
    if (data != null) {
      return data(this);
    }
    return orElse();
  }
}

abstract class _RecipesStateData implements RecipesState {
  const factory _RecipesStateData({required RecipesModel recipes}) =
      _$_RecipesStateData;

  RecipesModel get recipes => throw _privateConstructorUsedError;
  @JsonKey(ignore: true)
  _$RecipesStateDataCopyWith<_RecipesStateData> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class _$RecipesStateErrorCopyWith<$Res> {
  factory _$RecipesStateErrorCopyWith(
          _RecipesStateError value, $Res Function(_RecipesStateError) then) =
      __$RecipesStateErrorCopyWithImpl<$Res>;
  $Res call({String? error});
}

/// @nodoc
class __$RecipesStateErrorCopyWithImpl<$Res>
    extends _$RecipesStateCopyWithImpl<$Res>
    implements _$RecipesStateErrorCopyWith<$Res> {
  __$RecipesStateErrorCopyWithImpl(
      _RecipesStateError _value, $Res Function(_RecipesStateError) _then)
      : super(_value, (v) => _then(v as _RecipesStateError));

  @override
  _RecipesStateError get _value => super._value as _RecipesStateError;

  @override
  $Res call({
    Object? error = freezed,
  }) {
    return _then(_RecipesStateError(
      error == freezed
          ? _value.error
          : error // ignore: cast_nullable_to_non_nullable
              as String?,
    ));
  }
}

/// @nodoc

class _$_RecipesStateError implements _RecipesStateError {
  const _$_RecipesStateError([this.error]);

  @override
  final String? error;

  @override
  String toString() {
    return 'RecipesState.error(error: $error)';
  }

  @override
  bool operator ==(dynamic other) {
    return identical(this, other) ||
        (other is _RecipesStateError &&
            (identical(other.error, error) ||
                const DeepCollectionEquality().equals(other.error, error)));
  }

  @override
  int get hashCode =>
      runtimeType.hashCode ^ const DeepCollectionEquality().hash(error);

  @JsonKey(ignore: true)
  @override
  _$RecipesStateErrorCopyWith<_RecipesStateError> get copyWith =>
      __$RecipesStateErrorCopyWithImpl<_RecipesStateError>(this, _$identity);

  @override
  @optionalTypeArgs
  TResult when<TResult extends Object?>({
    required TResult Function() initial,
    required TResult Function() loading,
    required TResult Function(RecipesModel recipes) data,
    required TResult Function(String? error) error,
  }) {
    return error(this.error);
  }

  @override
  @optionalTypeArgs
  TResult maybeWhen<TResult extends Object?>({
    TResult Function()? initial,
    TResult Function()? loading,
    TResult Function(RecipesModel recipes)? data,
    TResult Function(String? error)? error,
    required TResult orElse(),
  }) {
    if (error != null) {
      return error(this.error);
    }
    return orElse();
  }

  @override
  @optionalTypeArgs
  TResult map<TResult extends Object?>({
    required TResult Function(_RecipesStateInitial value) initial,
    required TResult Function(_RecipesStateLoading value) loading,
    required TResult Function(_RecipesStateData value) data,
    required TResult Function(_RecipesStateError value) error,
  }) {
    return error(this);
  }

  @override
  @optionalTypeArgs
  TResult maybeMap<TResult extends Object?>({
    TResult Function(_RecipesStateInitial value)? initial,
    TResult Function(_RecipesStateLoading value)? loading,
    TResult Function(_RecipesStateData value)? data,
    TResult Function(_RecipesStateError value)? error,
    required TResult orElse(),
  }) {
    if (error != null) {
      return error(this);
    }
    return orElse();
  }
}

abstract class _RecipesStateError implements RecipesState {
  const factory _RecipesStateError([String? error]) = _$_RecipesStateError;

  String? get error => throw _privateConstructorUsedError;
  @JsonKey(ignore: true)
  _$RecipesStateErrorCopyWith<_RecipesStateError> get copyWith =>
      throw _privateConstructorUsedError;
}
